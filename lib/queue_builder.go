package lib

import "errors"

type QueueBuilder struct {
	//queue for working
	QueueName string

	//DeviceIDs for working, if not empty, will merge with queue
	DeviceIDs []string

	//logger
	server Server
}

func NewQueueBuilder(q string, d []string, server Server) (*QueueBuilder) {
	return &QueueBuilder{QueueName:q, DeviceIDs:d, server:server}
}

func (q *QueueBuilder) ToDeviceQueue(Capacity int) (*DeviceQueue, error) {
	queue := NewQueueByCapacity(Capacity, q.server)

	err := q.processData(queue)
	if err != nil {
		return nil, err
	}

	return queue, nil
}

// async mode device queue
func (q *QueueBuilder) AsyncToDeviceQueue(Capacity int) (*DeviceQueue, error) {
	queue := NewQueueByCapacity(Capacity, q.server)

	//async process data
	go q.processData(queue)

	return queue, nil
}

func (q *QueueBuilder) processData(queue *DeviceQueue) (error) {
	if q.QueueName != "" {
		q.server.GetEnv().GetLogger().Println("Init DeviceQueue data from QueueSource: "+q.QueueName)

		qs, err:=NewQueueSource(q.QueueName, *q.server.GetEnv().GetQueueSourceConfig())
		if err != nil {
			return errors.New("Error when NewQueueSource(): " + err.Error())
		}

		data, err:=qs.GetData()
		if err != nil {
			return errors.New("Error when qs.GetData(): " + err.Error())
		}
		err = queue.AppendDataSource(data)
		if err != nil {
			return errors.New("Error when queue.AppendDataSource(): " + err.Error())
		}
	}

	if q.DeviceIDs != nil && len(q.DeviceIDs)>0 {
		q.server.GetEnv().GetLogger().Println("Init DeviceQueue data from DeviceIDs parameter.")
		err := queue.AppendDataSource(q.DeviceIDs)
		if err != nil {
			return errors.New("Error when queue.AppendDataSource(): " + err.Error())
		}
	}

	//send pending
	queue.SetStatus(DEVICE_QUEUE_STATUS_PENDING)

	//close when finish, need to after add data or will finish without sending
	queue.EnableCloseAfterSended()
	return nil
}