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
)

const (
	WORKER_STATUS_SPARE = iota
	WORKER_STATUS_RUNNING

	WORKER_COMMAND_SEND = iota
	WORKER_COMMAND_STOP
)

type Worker struct {
	Client          *apns.Client
	Status          int
	WorkerID        int

	//worker lock
	Lock            sync.Mutex

	//push worker poll belong to
	Pool            *Pool

	//need to initialize
	PushChannel     chan *WorkerRequeset
	//WorkerID
	ResponseChannel chan *WorkerResponse
}

type WorkerRequeset struct {
	msg *apns.Notification

	cmd int
}

//Create a new request
func NewWorkerRequeset(msg *apns.Notification, cmd int) *WorkerRequeset {
	return &WorkerRequeset{msg:msg, cmd:cmd}
}

type WorkerResponse struct {
	Response *apns.Response

	Error    error
}

// create new worker
func NewWorker(env *EnvInfo) (*Worker, error) {
	cert, err := GetCerts(env.CertPath, env.CertPassword)
	if err != nil {
		return nil, err
	}

	client := apns.NewClient(cert).Production()
	worker := &Worker{Client:client, Status:WORKER_STATUS_SPARE, PushChannel:make(chan *WorkerRequeset), ResponseChannel:make(chan *WorkerResponse)}

	return worker, nil
}

// this is a goroutine run
func (p *Worker) Run() {
	env.GetLogger().Println(p.GetWorkerName() + " started, wait for push task...")
	for {
		select {
		// need transfer by copy
		case request := <-p.PushChannel:
			if request == nil || request.cmd == WORKER_COMMAND_STOP {
				env.GetLogger().Println(p.GetWorkerName() + " receive terminate channel signal, will quit.")
				break
			}else {
				//return response
				resp, err := p.Push(request.msg)
				p.ResponseChannel <- &WorkerResponse{Response:resp, Error:err}
			}
		}
	}

	//done wait group
	p.Pool.wg.Done()
}

// this a goroutine run
func (p *Worker) Subscribe(list *DeviceQueue, msg *apns.Notification) {
	for {
		msgLocal := *msg
		msgLocal.DeviceToken = <-list.Channel

		request := NewWorkerRequeset(&msgLocal, WORKER_COMMAND_SEND)
		p.PushChannel <- request

		//finish
		<-p.ResponseChannel
	}
}

func (p *Worker) Push(msgLocal *apns.Notification) (*apns.Response, error) {
	p.Lock.Lock()

	// working now
	p.Status = WORKER_STATUS_RUNNING

	env.GetLogger().Println(p.GetWorkerName() + " #start# to push for DeviceToken: " + msgLocal.DeviceToken)
	start := time.Now().UnixNano()

	resp, err := p.Client.Push(msgLocal)
	if err != nil {
		errMsg := p.GetWorkerName() + " Error while worker.Push():" + err.Error()
		env.GetLogger().Println(errMsg)
		p.Pool.GetFailLogger().Println(p.GetWorkerName() + " " + msgLocal.DeviceToken)
		return nil, errors.New(errMsg)
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
	p.Status = WORKER_STATUS_SPARE

	return resp, err
}

func (p *Worker) GetWorkerName() (string) {
	return "worker_" + strconv.Itoa(p.WorkerID)
}

func GetCerts(path, password string) (tls.Certificate, error) {
	cert, pemErr := certificate.FromP12File(path, password)
	if pemErr != nil {
		return tls.Certificate{}, errors.New("Cert Error:" + pemErr.Error())
	}

	return cert, nil
}