// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package lib

import (
	"sync"
	"os"
	"errors"
	"io/ioutil"

	"zooinit/log"
	"bytes"
	"strconv"
	"strings"
)

const (
	QUEUE_DEFAULT_CAPACITY = 50

	//init
	DEVICE_QUEUE_STATUS_INIT = "init"
	//ready for sending
	DEVICE_QUEUE_STATUS_PENDING = "pending"
	//suspend sending
	DEVICE_QUEUE_STATUS_SUSPEND = "suspend"
	//finish sending
	DEVICE_QUEUE_STATUS_FINISH = "finish"
)

type DeviceQueue struct {
	//channel is a synchronization
	Channel            chan string

	//发送文件位置
	Position           int

	data               []string
	//data locker
	lock               sync.Mutex

	status             string
	queueChangeChannel chan bool

	//if false can append queue after finish sending
	CloseAfterSended   bool

	//logger
	server             Server
}

func NewQueueByPool(p *Pool, server Server) (*DeviceQueue) {
	//2 times of pool capacity
	return NewQueueByCapacity(p.Config.Capacity * 2, server)
}

func NewQueueByCapacity(Capacity int, server Server) (*DeviceQueue) {
	//Capacity equal to pool
	chanCreate := make(chan string, Capacity)

	return &DeviceQueue{Channel:chanCreate, Position:0, status:DEVICE_QUEUE_STATUS_INIT, queueChangeChannel:make(chan bool, Capacity), CloseAfterSended:false, server:server}
}

func NewQueue(server Server) (*DeviceQueue) {
	return NewQueueByCapacity(QUEUE_DEFAULT_CAPACITY, server)
}

func NewQueueByServer(server Server) (*DeviceQueue) {
	//2 times of pool capacity
	return NewQueueByCapacity(server.GetEnv().GetPoolConfig().Capacity * 2, server)
}

//publish goroutine
//if status equal to init or suspend will block until data ready
func (q *DeviceQueue) Publish() {
	q.server.GetEnv().GetLogger().Println("DeviceQueue status is " + q.status + ", publish now...")

	for {
		if q.status == DEVICE_QUEUE_STATUS_INIT || q.status == DEVICE_QUEUE_STATUS_SUSPEND {
			for {
				q.server.GetEnv().GetLogger().Println("DeviceQueue status is " + q.status + ", will block q.queueChangeChannel...")

				//block
				<-q.queueChangeChannel

				if q.status != DEVICE_QUEUE_STATUS_INIT && q.status != DEVICE_QUEUE_STATUS_SUSPEND {
					q.server.GetEnv().GetLogger().Println("DeviceQueue status is " + q.status + ", will break wait for work.")
					//need to break loop
					break
				}
			}
		}

		//publish actual action
		q.sendToChannel()

		//finish work
		if q.status == DEVICE_QUEUE_STATUS_FINISH {
			break
		}
	}
}

func (q *DeviceQueue) sendToChannel() {
	// add a critical lock
	q.lock.Lock()
	defer q.lock.Unlock()

	//Pending need seding
	if q.status == DEVICE_QUEUE_STATUS_PENDING && q.Position < len(q.data) {
		q.Channel <- q.data[q.Position]
		q.Position++
	} else {
		if q.CloseAfterSended {
			//finish seding
			q.status = DEVICE_QUEUE_STATUS_FINISH
		}

		if q.status == DEVICE_QUEUE_STATUS_FINISH {
			//finish work
			close(q.Channel)
		}
	}
}

func (q *DeviceQueue) EnableCloseAfterSended() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.CloseAfterSended = true
}

func (q *DeviceQueue) DisableCloseAfterSended() {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.CloseAfterSended = false
}

func (q *DeviceQueue) TriggerChange() {
	q.queueChangeChannel <- true
}

// Status
// init->pending->finish(can goback to init)
//        ⬇️⬆️
//       suspend
func (q *DeviceQueue) SetStatus(status string) (bool, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.status == DEVICE_QUEUE_STATUS_INIT {
		if status != DEVICE_QUEUE_STATUS_PENDING {
			return false, errors.New("Not allowed to set status to " + status + ", NOW: " + q.status)
		} else {
			q.status = status
		}

	} else if q.status == DEVICE_QUEUE_STATUS_PENDING {
		if status == DEVICE_QUEUE_STATUS_SUSPEND || status == DEVICE_QUEUE_STATUS_FINISH {
			q.status = status

		} else {
			return false, errors.New("Not allowed to set status to " + status + ", NOW: " + q.status)
		}

	} else if q.status == DEVICE_QUEUE_STATUS_SUSPEND {
		if status != DEVICE_QUEUE_STATUS_PENDING {
			return false, errors.New("Not allowed to set status to " + status + ", NOW: " + q.status)
		} else {
			q.status = status
		}

	} else if q.status == DEVICE_QUEUE_STATUS_FINISH {
		if status != DEVICE_QUEUE_STATUS_INIT {
			return false, errors.New("Not allowed to set status to " + status + ", NOW: " + q.status)
		} else {
			q.status = status
			//rewind pos
			q.Position = 0
		}
	} else {
		return false, errors.New("Not support DeviceQueue status code.")
	}

	q.TriggerChange()
	return true, nil
}

func (q *DeviceQueue) ChangePosition(posNew int) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if posNew < len(q.data) && posNew >= 0 {
		q.Position = posNew
	}

	q.TriggerChange()
}

//publish goroutine
func (q *DeviceQueue) AppendFileDataSource(filename string) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	file, err := os.OpenFile(filename, os.O_RDONLY, log.DEFAULT_LOGFILE_MODE)
	if err != nil {
		return errors.New("DeviceQueue.AppendFileDataSource(): " + err.Error())
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return errors.New("DeviceQueue.AppendFileDataSource(): " + err.Error())
	}

	list := bytes.Split(content, []byte("\n"))

	for key, value := range list {
		//least string conversion
		err := q.appendInternalData(key, string(value))
		if err != nil {
			return err
		}
	}

	q.TriggerChange()

	return nil
}

//publish goroutine
func (q *DeviceQueue) AppendDataSource(list []string) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	for key, value := range list {
		err := q.appendInternalData(key, value)
		if err != nil {
			return err
		}
	}

	q.TriggerChange()

	return nil
}

func (q *DeviceQueue) appendInternalData(key int, value string) error {
	value = strings.Trim(value, "\n\r ")

	//TODO different platfrom device token length is different
	if len(value) == 64 {
		q.data = append(q.data, value)
	} else if len(value) == 0 {
		//may last line
		return nil
	} else {
		return errors.New("DeviceQueue.appendInternalData() error device token length: line " + strconv.Itoa(key + 1) + " -> " + value)
	}
	return nil
}

func (q *DeviceQueue) Len() int {
	return len(q.data)
}


