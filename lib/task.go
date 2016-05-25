package lib

import (
	"errors"
	"strconv"
	"sync"
)

const (
	TASK_QUEUE_MAX_WAITING = 100
	TASK_QUEUE_MAX_POOL = 5
)

type Task struct {
	// task device queue
	list    *DeviceQueue

	// sending message
	message MessageInterface
}

// task queue, cycle array
type TaskQueue struct {
	server          Server

	tasks           []*Task

	pool            []*Pool

	Lock            sync.Mutex

	readIndex       int
	writeIndex      int

	PushChannel     chan *Task
	ResponseChannel chan bool
}

func NewTaskQueue(server Server) *TaskQueue {
	return &TaskQueue{pool:make([]*Pool, TASK_QUEUE_MAX_POOL), tasks:make([]*Task, TASK_QUEUE_MAX_WAITING), server:server}
}

// add a new task
func (tq *TaskQueue)Add(list *DeviceQueue, msg MessageInterface) (int, error) {
	tq.Lock.Lock()
	defer tq.Lock.Unlock()

	if len(tq.tasks) >= TASK_QUEUE_MAX_WAITING {
		return 0, errors.New("Failed, Max task queue limit reached, limit: " + strconv.Itoa(TASK_QUEUE_MAX_WAITING))
	}

	task := &Task{list:list, message:msg}
	tq.tasks = append(tq.tasks, task)

	return len(tq.tasks), nil
}

// run task queue dispatch run
func (tq *TaskQueue) Run() {
	//initilize pools and pick one to run


}
