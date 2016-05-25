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

	// task queue locker
	Lock            sync.Mutex

	//Now Task Read Index
	readIndex       int
	//Now Task Write Index
	writeIndex      int

	PushChannel     chan *Task
	ResponseChannel chan bool
}

func NewTaskQueue(server Server) *TaskQueue {
	return &TaskQueue{pool:make([]*Pool, TASK_QUEUE_MAX_POOL), tasks:make([]*Task, TASK_QUEUE_MAX_WAITING), server:server}
}

func (tq *TaskQueue)nextID(index int) (int) {
	if (index > len(tq.tasks) - 1) {
		panic("TaskQueue op index out of bound.")
	}else if (index == len(tq.tasks) - 1) {
		return 0
	}else {
		return index + 1
	}
}

func (tq *TaskQueue)NextReadIndex() (int, error) {
	if tq.tasks[tq.nextID(tq.readIndex)] == nil {
		return 0, errors.New("Task Queue is empty now.")
	}else {
		return tq.nextID(tq.readIndex), nil
	}
}

func (tq *TaskQueue)NextWriteIndex() (int, error) {
	if tq.tasks[tq.writeIndex] == nil {
		// init state.
		return tq.writeIndex, nil
	}else if tq.tasks[tq.nextID(tq.writeIndex)] != nil {
		return 0, errors.New("Task Queue is full now, please wait...")
	}else {
		return tq.nextID(tq.writeIndex), nil
	}
}

// add a new task
func (tq *TaskQueue)Add(list *DeviceQueue, msg MessageInterface) (int, error) {
	tq.Lock.Lock()
	defer tq.Lock.Unlock()

	if list == nil {
		return 0, errors.New("Failed, invalid DeviceQueue.")
	}

	index, err := tq.NextWriteIndex()
	if err != nil {
		return 0, errors.New("Failed, " + err.Error() + ", limit: " + strconv.Itoa(TASK_QUEUE_MAX_WAITING))
	}

	task := &Task{list:list, message:msg}
	tq.tasks[index] = task

	//edit index
	tq.writeIndex = index

	pos := tq.writeIndex - tq.readIndex
	if pos < 0 {
		pos += len(tq.tasks)
	}

	return pos, nil
}

// pop now read task
func (tq *TaskQueue)Pop() (error) {
	tq.Lock.Lock()
	defer tq.Lock.Unlock()

	if tq.tasks[tq.readIndex] == nil {
		return errors.New("TaskQueue now is empty.")
	}
	tq.tasks[tq.readIndex] = nil

	//edit index
	index, err := tq.NextReadIndex()
	//empty not edit index
	if err == nil {
		tq.readIndex = index
	}

	return nil
}

// pop now read task
func (tq *TaskQueue)Read() (*Task, error) {
	if tq.tasks[tq.readIndex] == nil {
		return nil, errors.New("TaskQueue now is empty.")
	}else {
		return tq.tasks[tq.readIndex], nil
	}
}

// run task queue dispatch run
func (tq *TaskQueue) Run() {
	//initilize pools and pick one to run


}
