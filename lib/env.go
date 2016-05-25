// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"zooinit/cluster"
)

// This basic discovery service bootstrap env info
type EnvInfo interface {
	cluster.Env

	CreateWorker() (Worker, error)

	DestroyWorker(worker Worker) (error)

	GetPoolConfig() (*PoolConfig)
}
