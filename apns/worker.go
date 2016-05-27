// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package apns

import (
	"crypto/tls"
	"errors"
	"sync"
	"time"
	"strconv"

	apns "github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
	"github.com/sideshow/apns2/payload"

	"gopush/lib"
)

type Worker struct {
	Client          *apns.Client
	Status          int
	WorkerID        int

	//worker lock
	Lock            sync.Mutex

	//push worker poll belong to
	Pool            *lib.Pool

	//need to initialize
	PushChannel     chan *lib.WorkerRequeset
	//WorkerID
	ResponseChannel chan *lib.WorkerResponse
}


// create new worker
func NewWorker(env *EnvInfo) (*Worker, error) {
	cert, err := GetCerts(env.CertPath, env.CertPassword)
	if err != nil {
		return nil, err
	}

	client := apns.NewClient(cert).Production()
	worker := &Worker{Client:client, Status:lib.WORKER_STATUS_SPARE, PushChannel:make(chan *lib.WorkerRequeset), ResponseChannel:make(chan *lib.WorkerResponse)}

	return worker, nil
}

// this is a goroutine run
func (w *Worker) Run() {
	env.GetLogger().Println(w.GetWorkerName() + " started, wait for push task...")
	for {
		select {
		// need transfer by copy
		case request := <-w.PushChannel:
			if request == nil || request.Cmd == lib.WORKER_COMMAND_STOP {
				env.GetLogger().Println(w.GetWorkerName() + " receive terminate channel signal, will quit.")
				break
			}else {
				//return response
				resp := w.Push(request.Message, request.Device)
				w.ResponseChannel <- resp
			}
		}
	}
}

// this a goroutine run
func (w *Worker) Subscribe(task *lib.Task) {
	env.GetLogger().Println(w.GetWorkerName() + " started to Subscribe...")
	for {
		DeviceToken, more := <-task.GetList().Channel
		if more {
			request := lib.NewWorkerRequeset(task.GetMessage(), DeviceToken, lib.WORKER_COMMAND_SEND)
			w.PushChannel <- request

			//finish
			<-w.ResponseChannel
		}else {
			break
		}
	}
}

func (w *Worker) Push(msg lib.MessageInterface, Device string) (*lib.WorkerResponse) {
	w.Lock.Lock()

	msgLocal := &apns.Notification{}
	msgLocal.DeviceToken = Device
	msgLocal.ApnsID = msg.GetUuid()
	msgLocal.Priority = 10
	msgLocal.Topic = ""
	load := payload.NewPayload()

	load.Badge(1)
	load.AlertTitle(msg.GetTitle())
	load.AlertBody(msg.GetBody())
	//Done push Turn to specific page machanism, addon field
	if msg.GetCustom() != nil {
		haimiPayloadKey := "payload"
		if haimiPayload, ok := msg.GetCustom()[haimiPayloadKey]; ok {
			if len(haimiPayload) > 0 {
				load.Custom(haimiPayloadKey, msg.GetCustom()[haimiPayloadKey])
			}
		}
	}

	load.Sound(msg.GetSound())
	msgLocal.Payload = load

	// working now
	w.Status = lib.WORKER_STATUS_RUNNING

	env.GetLogger().Println(w.GetWorkerName() + " #start# to push for DeviceToken: " + msgLocal.DeviceToken)
	start := time.Now().UnixNano()

	resp, err := w.Client.Push(msgLocal)
	if err != nil {
		errMsg := w.GetWorkerName() + " Error while worker.Push():" + err.Error()
		env.GetLogger().Println(errMsg)
		w.Pool.GetFailLogger().Println(w.GetWorkerName() + " " + msgLocal.DeviceToken)
		return &lib.WorkerResponse{Response:nil, Error:errors.New(errMsg)}
	}

	//in us
	timeSpent := (time.Now().UnixNano() - start) / 1000
	//success
	if resp.Sent() {
		env.GetLogger().Println(w.GetWorkerName() + " sent #success#: " + msgLocal.DeviceToken + " -> " + resp.ApnsID)
		w.Pool.GetOKLogger().Println(w.GetWorkerName() + " " + msgLocal.DeviceToken + " -> " + resp.ApnsID + " -> " + strconv.Itoa(int(timeSpent)) + "us")
	}else {
		env.GetLogger().Println(w.GetWorkerName() + " sent #faild#: " + msgLocal.DeviceToken + " -> " + resp.Reason)
		w.Pool.GetFailLogger().Println(w.GetWorkerName() + " " + msgLocal.DeviceToken + " -> " + strconv.Itoa(resp.StatusCode) + " -> " + resp.Reason + " -> " + strconv.Itoa(int(timeSpent)) + "us -> " + resp.Timestamp.Format(time.RFC3339))
	}

	w.Lock.Unlock()
	w.Status = lib.WORKER_STATUS_SPARE

	return &lib.WorkerResponse{Response:resp, Error:err}
}

func (w *Worker) GetWorkerName() (string) {
	return "pool_" + strconv.Itoa(w.Pool.PoolID) + "_worker_" + strconv.Itoa(w.WorkerID)
}

func (w *Worker) SetWorkerID(id int) (bool) {
	w.WorkerID = id
	return true
}

func (w *Worker) SetPool(pool *lib.Pool) (bool) {
	w.Pool = pool
	return true
}

func (w *Worker) Destroy() (error) {

	return nil
}

func GetCerts(path, password string) (tls.Certificate, error) {
	cert, pemErr := certificate.FromP12File(path, password)
	if pemErr != nil {
		return tls.Certificate{}, errors.New("Cert Error:" + pemErr.Error())
	}

	return cert, nil
}

