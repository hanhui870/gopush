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

func (q*QueueBuilder) ToDeviceQueue(Capacity int) (*DeviceQueue, error) {
	queue := NewQueueByCapacity(Capacity)

	if q.QueueName != "" {
		err := queue.AppendFileDataSource("runtime/data/" + q.QueueName + ".txt")
		if err != nil {
			return nil, err
		}
	}

	if q.DeviceIDs != nil {
		err := queue.AppendDataSource(q.DeviceIDs)
		if err != nil {
			return nil, err
		}
	}

	return queue, nil
}