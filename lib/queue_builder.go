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
	//use default
	if q.DeviceIDs==nil && q.QueueName == "" {
		q.QueueName=q.server.GetEnv().GetQueueSourceConfig().Value
	}

	if q.QueueName != "" {
		q.server.GetEnv().GetLogger().Println("Init DeviceQueue data from QueueSource: "+q.QueueName)

		qs, err:=NewQueueSource(q.QueueName, *q.server.GetEnv().GetQueueSourceConfig())
		if err != nil {
			msg:="Error when NewQueueSource(): " + err.Error()
			q.server.GetEnv().GetLogger().Println(msg)
			return errors.New(msg)
		}

		data, err:=qs.GetData()
		if err != nil {
			msg:="Error when qs.GetData(): " + err.Error()
			q.server.GetEnv().GetLogger().Println(msg)
			return errors.New(msg)
		}
		err = queue.AppendDataSource(data)
		if err != nil {
			msg:="Error when queue.AppendDataSource(): " + err.Error()
			q.server.GetEnv().GetLogger().Println(msg)
			return errors.New(msg)
		}
	}

	if q.DeviceIDs != nil && len(q.DeviceIDs)>0 {
		q.server.GetEnv().GetLogger().Println("Init DeviceQueue data from DeviceIDs parameter.")
		err := queue.AppendDataSource(q.DeviceIDs)
		if err != nil {
			msg:="Error when queue.AppendDataSource(): " + err.Error()
			q.server.GetEnv().GetLogger().Println(msg)
			return errors.New(msg)
		}
	}

	//send pending
	queue.SetStatus(DEVICE_QUEUE_STATUS_PENDING)

	//close when finish, need to after add data or will finish without sending
	queue.EnableCloseAfterSended()
	return nil
}