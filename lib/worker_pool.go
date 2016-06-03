// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"errors"
	"sync"
)

// A global worker env related management
// worker identify use uuid
type WorkerPool struct {
	//string is uuid field
	workers map[string]Worker

	//Env related workers
	env     EnvInfo

	lock    sync.Mutex
}

func NewWorkerPool(env EnvInfo) (*WorkerPool, error) {
	// Max workers controlled by poolNum*pool.Capacity
	// make map: an optional capacity hint, The initial capacity does not bound its size: maps grow to accommodate the number of items stored in them, with the exception of nil maps.
	return &WorkerPool{env:env, workers:make(map[string]Worker, POOL_DEFAULT_CAPACITY * TASK_QUEUE_MAX_POOL)}, nil
}

func (wp *WorkerPool) CreateWorker() (Worker, error) {
	wp.lock.Lock()
	defer wp.lock.Unlock()

	var worker Worker
	for _, wkTmp := range wp.workers {
		if wkTmp.GetStatus() == WORKER_STATUS_SPARE {
			worker = wkTmp
		}
	}

	if worker == nil {
		var err error
		worker, err = wp.env.CreateWorker()
		if err != nil {
			return nil, err
		}

		wp.workers[worker.GetUUID()] = worker
	}

	return worker, nil
}

//reset Worker ownership
func (wp *WorkerPool) HarvestWorker(worker Worker) (error) {
	(*worker.GetLockPtr()).Lock()
	defer (*worker.GetLockPtr()).Unlock()

	if worker.GetStatus() != WORKER_STATUS_SPARE {
		return errors.New("Error when WorkerPool.HarvestWorker(): worker.GetStatus!=WORKER_STATUS_SPARE")
	}

	if wk, ok := wp.workers[worker.GetUUID()]; ok {
		wk.SetPool(nil)
		wk.SetWorkerID(0)
	}

	return nil
}



