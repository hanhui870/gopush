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
	server            Server

	tasks             []*Task
	taskChangeChannel chan bool

	pools             []*Pool
	poolFinishChannel chan int

	// task queue locker
	Lock              sync.Mutex

	//Now Task Read Index
	readIndex         int
	//Now Task Write Index
	writeIndex        int

	wg                sync.WaitGroup
}

func NewTaskQueue(server Server) *TaskQueue {
	//PublishChannel no buffer
	return &TaskQueue{pools:make([]*Pool, TASK_QUEUE_MAX_POOL), tasks:make([]*Task, TASK_QUEUE_MAX_WAITING),
		taskChangeChannel:make(chan bool, TASK_QUEUE_MAX_WAITING), poolFinishChannel:make(chan int, TASK_QUEUE_MAX_WAITING), server:server}
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

	tq.taskChangeChannel <- true

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

//publish goroutine
func (tq *TaskQueue) getSparePool() (*Pool) {
	for _, pool := range tq.pools {
		if pool != nil && pool.Status == POOL_STATUS_SPARE {
			return pool
		}
	}

	return nil
}

//publish goroutine
func (tq *TaskQueue) publish() {
	for {
		task, err := tq.Read()
		if err != nil {
			// empty
			<-tq.taskChangeChannel
		}else {
			//select pool or create
			//spare pool -> create pool -> wait
			var poolSelected *Pool
			pool := tq.getSparePool()
			if pool != nil {
				poolSelected = pool
			}else {
				for iter, pool := range tq.pools {
					if (pool == nil) {
						//need a clone's pointer
						config := &(*tq.server.GetEnv().GetPoolConfig())
						config.SetSizeByQueueLength(task.list.Len())
						pool, err = NewPoolByConfig(config, tq.server.GetEnv())
						if err != nil {
							tq.server.GetEnv().GetLogger().Println("Create pool failed:" + err.Error())
						}else {
							//update poolid
							pool.PoolID = iter

							tq.pools[iter] = pool
							poolSelected = pool
							break
						}
					}
				}
			}

			if poolSelected != nil {
				go func() {
					//triger sending
					poolSelected.Send(task, tq.poolFinishChannel)
				}()

				//pop task when started, or will resend
				//TODO if send failed
				tq.Pop()

			}else {
				//read notify
				<-tq.poolFinishChannel
			}
		}
	}
}

// run task queue dispatch run
func (tq *TaskQueue) Run() {
	//initilize pools and pick one to run
	tq.wg.Add(1)
	go func() {
		tq.publish()

		tq.wg.Done()
	}()

	tq.wg.Wait()
}

func (t *Task) GetList() *DeviceQueue {
	return t.list
}

func (t *Task) GetMessage() MessageInterface {
	return t.message
}
