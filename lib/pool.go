// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"sync"
	"log"
	"errors"
//"time"

	loglocal "zooinit/log"
	"strconv"
)

const (
	POOL_STATUS_SPARE = iota
	POOL_STATUS_RUNNING

	POOL_DEFAULT_SIZE = 5
	POOL_DEFAULT_CAPACITY = 500
//test multi workers can set to 2
	POOL_DEFAULT_MINISPARE = 2
	POOL_DEFAULT_MAXSPARE = 50
)

//pool automatic resize if needed
//TODO can use workers globally for connection saving. @see initWorkers()
type Pool struct {
	//worker pool
	Workers    []Worker

	//Pool status
	Status     int
	PoolID     int

	Config     *PoolConfig

	//pool lock
	Lock       sync.Mutex

	//worker wg
	wg         sync.WaitGroup

	OKLogger   *loglocal.BufferedFileLogger
	FailLogger *loglocal.BufferedFileLogger

	Env        EnvInfo

	//sending wg
	sendWg     sync.WaitGroup

	//Every related to a Task, which can be changed every run.
	task       *Task
}

type PoolConfig struct {
	//worker poll size
	//MiniSpare <= now <= Capacity
	Size      int

	//pool capacity
	Capacity  int

	//mini spare worker
	MiniSpare int

	//max spare worker
	MaxSpare  int
}


// create a new worker pool
func NewPool(Size, Capacity, MiniSpare, MaxSpare int, Env EnvInfo) (*Pool, error) {
	config, err := NewPoolConfig(Size, Capacity, MiniSpare, MaxSpare)
	if err != nil {
		return nil, err
	}

	return NewPoolByConfig(config, Env)
}

func NewPoolByConfig(config *PoolConfig, Env EnvInfo) (*Pool, error) {
	pool := &Pool{Config:config}
	pool.Env = Env

	err:=pool.initWorkers(pool.Config.Size)
	if err!=nil {
		return nil, errors.New("Error when NewPoolByConfig():" + err.Error())
	}

	//run when created.
	go pool.Run()

	return pool, nil
}

func (p *Pool) initWorkers(NewCount int) error {
	oldWorkers := p.Workers
	// initWorkers Need to use new count
	workers := make([]Worker, NewCount, p.Config.Capacity)

	for iter, _ := range workers {
		var worker Worker
		var err error

		//fetch old and reuse worker, length compare
		if oldWorkers != nil && len(oldWorkers) <= iter + 1 {
			worker = oldWorkers[iter]
		}else {
			worker, err = p.Env.CreateWorker()
			if err != nil {
				return err
			}
		}


		//不能用append, 会增长数组.
		//workers=append(workers, worker)
		//start from 0
		worker.SetWorkerID(iter)
		worker.SetPool(p)

		workers[iter] = worker
	}

	// need to destroy old workers
	// TODO This workers can be reusable, but have to related to Env for multi certs.
	if len(oldWorkers) > len(workers) {
		iter := len(workers)//no need minus 1
		for ; iter < len(oldWorkers); iter++ {
			//trigger action
			oldWorkers[iter].Destroy()
		}
	}

	// edit new count
	p.Config.Size=NewCount
	p.Workers = workers

	return nil
}

func (p *Pool) Run() {
	if len(p.Workers) != p.Config.Size {
		p.Env.GetLogger().Fatalln("Found exception of pool: len(p.Workers)!=p.Size: ", len(p.Workers), p.Config.Size)
	}

	for _, worker := range p.Workers {
		p.wg.Add(1)
		//env.GetLogger().Println(worker.GetWorkerName()+" ...")

		//fixed: Here has an error Mode if run in anonymous func, worker started is not in expected mode
		go func(worker Worker) {
			worker.Run()
			p.wg.Done()
		}(worker)
	}

	// wait for all done worker.Run() / worker.Subscribe()
	p.wg.Wait()
}

// finish, taskqueue's finish channel
func (p *Pool) Send(task *Task, finish chan int) {
	p.Lock.Lock()
	p.Status = POOL_STATUS_RUNNING

	con, err := task.message.MarshalJSON()
	if err != nil {
		p.Env.GetLogger().Println("msg.MarshalJSON() found error:", err)
		return
	}
	p.Env.GetLogger().Println("Receive new push task: " + string(con))

	p.sendWg.Add(1)
	// Queue data publish
	go func() {
		task.list.Publish()

		p.sendWg.Done()
	}()

	for _, worker := range p.Workers {
		p.sendWg.Add(1)
		go func(worker Worker) {
			worker.Subscribe(task)
			p.sendWg.Done()
		}(worker)
	}

	p.sendWg.Wait()

	//test, pools iter
	//time.Sleep(5*time.Second)

	p.Lock.Unlock()
	p.Status = POOL_STATUS_SPARE
	finish <- p.PoolID
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

func (p *Pool) GetTask() (*Task) {
	return p.task
}

//Resize pool worker pools
func (p *Pool) Resize(size int) (error) {
	p.Lock.Lock()
	defer p.Lock.Unlock()

	//need a clone's pointer
	pCfg:=*p.Config
	pCfgNew := &pCfg

	pCfgNew.SetSizeByQueueLength(size)

	//new<old
	if pCfgNew.Size < p.Config.Size {
		return p.harvest(pCfgNew.Size)

	}else if pCfgNew.Size > p.Config.Size {
		return p.expand(pCfgNew.Size)
	}

	return nil
}

//can add when running
func (p *Pool) expand(size int) (error) {
	return p.initWorkers(size)
}

func (p *Pool) GetPoolName() string {
	return "pool_"+strconv.Itoa(p.PoolID)
}

//can not harvest when running
func (p *Pool) harvest(size int) (error) {
	if p.Status == POOL_STATUS_RUNNING {
		return errors.New("Pool can't add worker when running.")
	}

	return p.initWorkers(size)
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

func (pc *PoolConfig) SetSizeByQueueLength(length int) {
	if length <= 10 {
		pc.Size = pc.MiniSpare
	}else if length <= 100 {
		pc.Size = pc.MiniSpare * 5
	}else if length <= 1000 {
		pc.Size = pc.MiniSpare * 10
	}else if length <= 10000 {
		pc.Size = pc.MiniSpare * 50
	}else if length <= 100000 {
		pc.Size = pc.MiniSpare * 100
	}else if length <= 500000 {
		pc.Size = pc.MiniSpare * 150
	}else {
		pc.Size = pc.Capacity
	}

	if pc.Size > pc.Capacity {
		pc.Size = pc.Capacity
	}
}