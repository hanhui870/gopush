// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"sync"
	"log"
	"errors"

	loglocal "zooinit/log"
)

type Pool struct {
	//worker pool
	Workers       []Worker
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

	//worker wg
	wg            sync.WaitGroup

	OKLogger      *loglocal.BufferedFileLogger
	FailLogger    *loglocal.BufferedFileLogger

	Env           EnvInfo

	//sending wg
	sendWg        sync.WaitGroup

	task          *TaskQueue
}

// create a new worker pool
func NewPool(Size, Capacity, MiniSpare, MaxSpare int, Env EnvInfo) (*Pool, error) {
	workers := make([]Worker, Size, Capacity)
	WorkerIDIndex := 0

	pool := &Pool{Size:Size, Capacity:Capacity, MaxSpare:MaxSpare, MiniSpare:MiniSpare}

	for iter, _ := range workers {
		worker, err := Env.CreateWorker()
		if err != nil {
			return nil, err
		}

		//不能用append, 会增长数组.
		//workers=append(workers, worker)
		//start from 0
		worker.SetWorkerID(iter)
		WorkerIDIndex = iter

		worker.SetPool(pool)

		workers[iter] = worker
	}

	//WorkerIDIndex is last new one
	pool.WorkerIDIndex = WorkerIDIndex + 1
	pool.Workers = workers
	pool.Env = Env

	return pool, nil
}

func (p *Pool) FetchASpareWork() Worker {
	return nil
}

func (p *Pool) Run() {
	if len(p.Workers) != p.Size {
		p.Env.GetLogger().Fatalln("Found exception of pool: len(p.Workers)!=p.Size: ", len(p.Workers), p.Size)
	}

	// start up worker
	for _, worker := range p.Workers {
		p.wg.Add(1)
		//env.GetLogger().Println(worker.GetWorkerName()+" ...")

		//TODO Here has an error Mode if run in anonymous func, worker started is not in expected mode
		go func() {
			worker.Run()
			p.wg.Done()
		}()

	}

	// wait for all done worker.Run() / worker.Subscribe()
	p.wg.Wait()
}

func (p *Pool) Send(list *DeviceQueue, msg MessageInterface) {
	con, err := msg.MarshalJSON()
	if err != nil {
		p.Env.GetLogger().Println("msg.MarshalJSON() found error:", err)
		return
	}
	p.Env.GetLogger().Println("Receive new push task: " + string(con))

	p.sendWg.Add(1)
	// Queue data publish
	go func() {
		list.Publish()

		p.sendWg.Done()
	}()

	for _, worker := range p.Workers {
		p.sendWg.Add(1)
		go func() {
			worker.Subscribe(list, msg)
			p.sendWg.Done()
		}()
	}

	p.sendWg.Wait()
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

func (p *Pool) GetTaskQueue() (*TaskQueue) {
	return p.task
}

func (p *Pool) getInternalLogger(logtype string) (*loglocal.BufferedFileLogger) {
	filename := loglocal.GenerateFileLogPathName(p.Env.GetLogPath(), "pool_" + logtype)
	file, err := loglocal.NewFileLog(filename)
	if err != nil {
		log.Fatalln(err)
	}

	logger := log.New(file, "", log.Ldate | log.Ltime | log.Lmicroseconds) // add time for stat
	return loglocal.GetBufferedFileLogger(file, logger)
}

type PoolConfig struct {
	Size      int
	Capacity  int
	MiniSpare int
	MaxSpare  int
}

func NewPoolConfig(Size, Capacity, MiniSpare, MaxSpare int) (*PoolConfig, error) {
	if Size <= 0 || Capacity <= 0 || MiniSpare <= 0 || MaxSpare <= 0 {
		return nil, errors.New("All Size, Capacity, MiniSpare, MaxSpare parameters must all >0")
	}
	if Size < MiniSpare {
		Size = MiniSpare
		return nil, errors.New("Size<MiniSpare, will change to equal to MiniSpare")
	}
	if Size > Capacity {
		Size = Capacity
		return nil, errors.New("Size>Capacity, will change to equal to Capacity")
	}
	if MiniSpare > MaxSpare {
		return nil, errors.New("MiniSpare must <= MaxSpare")
	}
	if Size > Capacity || MiniSpare > Capacity || MaxSpare >= Capacity {
		return nil, errors.New("Capacity must be the greatest parameter within Size, Capacity, MiniSpare, MaxSpare")
	}

	return &PoolConfig{Size:Size, Capacity:Capacity, MiniSpare:MiniSpare, MaxSpare:MaxSpare}, nil
}