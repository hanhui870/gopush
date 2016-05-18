// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import ()

const (
	WORKER_STATUS_SPARE = iota
	WORKER_STATUS_RUNNING

	WORKER_COMMAND_SEND = iota
	WORKER_COMMAND_STOP
)

type Worker interface {
	Run()

	// Subscribe device goroutine run
	Subscribe(list *DeviceQueue, msg Message)

	// Push a message
	Push(msgLocal Message) (*WorkerResponse)

	// Fetch a worker identify
	GetWorkerName() (string)

	// Set worker id
	SetWorkerID(id int) (bool)

	// Set worker's belonging pool
	SetPool(pool *Pool) (bool)
}

type WorkerRequeset struct {
	Message interface{}

	Cmd     int
}

//Create a new request
func NewWorkerRequeset(msg interface{}, cmd int) *WorkerRequeset {
	return &WorkerRequeset{Message:msg, Cmd:cmd}
}

type WorkerResponse struct {
	Response interface{}

	Error    error
}