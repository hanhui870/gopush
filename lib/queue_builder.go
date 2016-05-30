package lib

type QueueBuilder struct {
	//queue for working
	QueueName string

	//DeviceIDs for working, if not empty, will merge with queue
	DeviceIDs []string
}

func NewQueueBuilder(q string, d []string) (*QueueBuilder) {
	return &QueueBuilder{QueueName:q, DeviceIDs:d}
}

func (q *QueueBuilder) ToDeviceQueue(Capacity int) (*DeviceQueue, error) {
	queue := NewQueueByCapacity(Capacity)

	err := q.processData(queue)
	if err != nil {
		return nil, err
	}

	return queue, nil
}

// async mode device queue
func (q *QueueBuilder) AsyncToDeviceQueue(Capacity int) (*DeviceQueue, error) {
	queue := NewQueueByCapacity(Capacity)

	//async process data
	go q.processData(queue)

	return queue, nil
}

func (q *QueueBuilder) processData(queue *DeviceQueue) (error) {
	if q.QueueName != "" {
		err := queue.AppendFileDataSource("runtime/data/" + q.QueueName + ".txt")
		if err != nil {
			return err
		}
	}

	if q.DeviceIDs != nil {
		err := queue.AppendDataSource(q.DeviceIDs)
		if err != nil {
			return err
		}
	}

	queue.SetStatus(DEVICE_QUEUE_STATUS_SUSPEND)

	return nil
}