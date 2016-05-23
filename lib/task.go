package lib

import (
	"errors"
	"strconv"
	"sync"
)

const (
	TASK_QUEUE_MAX_WAITING = 100
)

type Task struct {
	// task device queue
	list    *DeviceQueue

	// sending message
	message MessageInterface
}

// task queue, cycle array
type TaskQueue struct {
	tasks      []*Task

	pool       *Pool

	Lock       sync.Mutex

	readIndex  int
	writeIndex int
}

func NewTaskQueue(pool *Pool) *TaskQueue {
	return &TaskQueue{pool:pool, tasks:make([]*Task, TASK_QUEUE_MAX_WAITING)}
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

	// notify
	tq.pool.taskQueueChannel <- true

	return len(tq.tasks), nil
}