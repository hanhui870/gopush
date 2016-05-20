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
)

type DeviceQueue struct {
	//channel is a synchronization
	Channel     chan string

	//发送文件位置
	Position    int

	data        []string
	//data locker
	lock        sync.Mutex
	//data change channel
	dataChannel chan bool
}

func NewQueueByPool(p *Pool) (*DeviceQueue) {
	//2 times of pool capacity
	return NewQueueByCapacity(p.Capacity * 2)
}

func NewQueueByCapacity(Capacity int) (*DeviceQueue) {
	//Capacity equal to pool
	chanCreate := make(chan string, Capacity)

	return &DeviceQueue{Channel:chanCreate, Position:0, dataChannel:make(chan bool, 10)}
}

func NewQueue() (*DeviceQueue) {
	return NewQueueByCapacity(QUEUE_DEFAULT_CAPACITY)
}

//publish goroutine
func (q *DeviceQueue) Publish() {
	for {
		if q.Position < len(q.data) {
			q.Channel <- q.data[q.Position]
			q.Position++
		}else {
			//wait data change
			<-q.dataChannel
		}
	}
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

	//notify data change
	q.dataChannel <- true
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

	//notify data change
	q.dataChannel <- true
	return nil
}

func (q *DeviceQueue) appendInternalData(key int, value string) error {
	value = strings.Trim(value, "\n\r ")
	if len(value) == 64 {
		q.data = append(q.data, value)
	}else {
		return errors.New("DeviceQueue.appendInternalData() error device token length: line " + strconv.Itoa(key + 1) + " -> " + value)
	}
	return nil
}

func (q *DeviceQueue) Len() int {
	return len(q.data)
}


