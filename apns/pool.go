// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package apns

import (
	"sync"
	"log"

	apns "github.com/sideshow/apns2"

	loglocal "zooinit/log"
	"strconv"
)

type Pool struct {
	//worker pool
	Workers       []*Worker
	WorkerIDIndex int

	//worker poll size
	//MiniSpare <= now <= Capacity
	Size          int

	//pool capacity
	Capacity      int

	//mini spare worker
	MiniSpare     int

	//max spare worker
	MaxSpare      int

	//pool lock
	Lock          sync.Mutex

	wg            sync.WaitGroup

	OKLogger      *loglocal.BufferedFileLogger
	FailLogger    *loglocal.BufferedFileLogger
}

// create a new worker pool
func NewPool(Size, Capacity, MiniSpare, MaxSpare int) (*Pool, error) {
	workers := make([]*Worker, Size, Capacity)
	WorkerIDIndex := 0

	pool := &Pool{Size:Size, Capacity:Capacity, MaxSpare:MaxSpare, MiniSpare:MiniSpare}

	for iter, _ := range workers {
		worker, err := NewWorker(env)
		if err != nil {
			return nil, err
		}

		//不能用append, 会增长数组.
		//workers=append(workers, worker)
		//start from 0
		worker.WorkerID = iter
		WorkerIDIndex = iter

		worker.Pool = pool

		workers[iter] = worker
	}

	//WorkerIDIndex is last new one
	pool.WorkerIDIndex = WorkerIDIndex + 1
	pool.Workers = workers

	return pool, nil
}

func (p *Pool) FetchASpareWork() *Worker {
	return nil
}

func (p *Pool) Run(list *DeviceQueue, msg *apns.Notification) {
	con, err := msg.MarshalJSON()
	if err != nil {
		env.GetLogger().Println("msg.MarshalJSON() found error:", err)
		return
	}
	env.GetLogger().Println("Receive new push task: " + string(con))

	if len(p.Workers) != p.Size {
		env.GetLogger().Fatalln("Found exception of pool: len(p.Workers)!=p.Size: ", len(p.Workers), p.Size)
	}

	// start up worker
	for _, worker := range p.Workers {
		p.wg.Add(1)
		//env.GetLogger().Println(worker.GetWorkerName()+" ...")

		//TODO Here has an error Mode if run in anonymous func, worker started is not in expected mode
		go worker.Run()

		//goroutine wg.done in worker.Run()
	}

	//p.Test(list, msg)
	p.Send(list, msg)

	// wait for all done worker.Run() / worker.Subscribe()
	p.wg.Wait()
}

// test code
func (p *Pool) Test(list *DeviceQueue, msg *apns.Notification) {
	//need to be copy
	msgLocal := *msg
	msgLocal.DeviceToken = "3523544012e5491b3fe8cf6627eddd123d6aa4191fbebf371191a3ce7d4c02ac"
	//failed test DeviceTokenNotForTopic
	//msgLocal.DeviceToken="fc9f5a80b3338e0259d235e3b7cef12d12137d888d8b7901005e258df5f1a863"

	request := NewWorkerRequeset(&msgLocal, WORKER_COMMAND_SEND)
	p.Workers[0].PushChannel <- request
	resp := <-p.Workers[0].ResponseChannel
	env.GetLogger().Println("Receive new push result: ", *resp)

	for key, value := range list.data {
		env.GetLogger().Println(strconv.Itoa(key) + " -> " + value)
	}
}

func (p *Pool) Send(list *DeviceQueue, msg *apns.Notification) {
	// Queue data publish
	go list.Publish()

	for _, worker := range p.Workers {

		p.wg.Add(1)
		go worker.Subscribe(list, msg)

		//goroutine wg.done in worker.Subscribe()
	}
}

func (p *Pool) GetOKLogger() (*loglocal.BufferedFileLogger) {
	if p.OKLogger == nil {
		p.OKLogger = p.getInternalLogger("ok")
	}

	return p.OKLogger
}

func (p *Pool) GetFailLogger() (*loglocal.BufferedFileLogger) {
	if p.FailLogger == nil {
		p.FailLogger = p.getInternalLogger("fail")
	}

	return p.FailLogger
}

func (p *Pool) getInternalLogger(logtype string) (*loglocal.BufferedFileLogger) {
	filename := loglocal.GenerateFileLogPathName(env.LogPath, "pool_" + logtype)
	file, err := loglocal.NewFileLog(filename)
	if err != nil {
		log.Fatalln(err)
	}

	logger := log.New(file, "", log.Ldate | log.Ltime | log.Lmicroseconds) // add time for stat
	return loglocal.GetBufferedFileLogger(file, logger)
}