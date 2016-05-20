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
	"github.com/twinj/uuid"

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
func (p *Worker) Run() {
	env.GetLogger().Println(p.GetWorkerName() + " started, wait for push task...")
	for {
		select {
		// need transfer by copy
		case request := <-p.PushChannel:
			if request == nil || request.Cmd == lib.WORKER_COMMAND_STOP {
				env.GetLogger().Println(p.GetWorkerName() + " receive terminate channel signal, will quit.")
				break
			}else {
				//return response
				resp := p.Push(request.Message, request.Device)
				p.ResponseChannel <- resp
			}
		}
	}
}

// this a goroutine run
func (p *Worker) Subscribe(list *lib.DeviceQueue, msg lib.MessageInterface) {
	for {
		DeviceToken := <-list.Channel

		request := lib.NewWorkerRequeset(msg, DeviceToken, lib.WORKER_COMMAND_SEND)
		p.PushChannel <- request

		//finish
		<-p.ResponseChannel
	}
}

func (p *Worker) Push(msg lib.MessageInterface, Device string) (*lib.WorkerResponse) {
	p.Lock.Lock()

	msgLocal := &apns.Notification{}
	msgLocal.DeviceToken = Device
	msgLocal.ApnsID = uuid.NewV1().String()
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
				load.Custom(haimiPayloadKey, "haimi-590")
			}
		}
	}

	load.Sound(msg.GetSound())
	msgLocal.Payload = load

	// working now
	p.Status = lib.WORKER_STATUS_RUNNING

	env.GetLogger().Println(p.GetWorkerName() + " #start# to push for DeviceToken: " + msgLocal.DeviceToken)
	start := time.Now().UnixNano()

	resp, err := p.Client.Push(msgLocal)
	if err != nil {
		errMsg := p.GetWorkerName() + " Error while worker.Push():" + err.Error()
		env.GetLogger().Println(errMsg)
		p.Pool.GetFailLogger().Println(p.GetWorkerName() + " " + msgLocal.DeviceToken)
		return &lib.WorkerResponse{Response:nil, Error:errors.New(errMsg)}
	}

	//in us
	timeSpent := (time.Now().UnixNano() - start) / 1000
	//success
	if resp.Sent() {
		env.GetLogger().Println(p.GetWorkerName() + " sent #success#: " + msgLocal.DeviceToken + " -> " + resp.ApnsID)
		p.Pool.GetOKLogger().Println(p.GetWorkerName() + " " + msgLocal.DeviceToken + " -> " + resp.ApnsID + " -> " + strconv.Itoa(int(timeSpent)) + "us")
	}else {
		env.GetLogger().Println(p.GetWorkerName() + " sent #faild#: " + msgLocal.DeviceToken + " -> " + resp.Reason)
		p.Pool.GetFailLogger().Println(p.GetWorkerName() + " " + msgLocal.DeviceToken + " -> " + strconv.Itoa(resp.StatusCode) + " -> " + resp.Reason + " -> " + strconv.Itoa(int(timeSpent)) + "us -> " + resp.Timestamp.Format(time.RFC3339))
	}

	p.Lock.Unlock()
	p.Status = lib.WORKER_STATUS_SPARE

	return &lib.WorkerResponse{Response:resp, Error:err}
}

func (p *Worker) GetWorkerName() (string) {
	return "worker_" + strconv.Itoa(p.WorkerID)
}

func (p *Worker) SetWorkerID(id int) (bool) {
	p.WorkerID = id
	return true
}

func (p *Worker) SetPool(pool *lib.Pool) (bool) {
	p.Pool = pool
	return true
}

func GetCerts(path, password string) (tls.Certificate, error) {
	cert, pemErr := certificate.FromP12File(path, password)
	if pemErr != nil {
		return tls.Certificate{}, errors.New("Cert Error:" + pemErr.Error())
	}

	return cert, nil
}