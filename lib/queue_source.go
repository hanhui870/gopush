package lib

import (
	"errors"
	"strings"
)

const (
	QUEUE_SOURCE_METHOD_API = "api"
	QUEUE_SOURCE_METHOD_FILE = "file"
	QUEUE_SOURCE_METHOD_MYSQL = "mysql"
)

type QueueSource struct {
	config *QueueSourceConfig
}

//queue data source fetch, SQL or restful API
type QueueSourceConfig struct {
	Method   string
	//sql dsn config
	MysqlDsn string
	//Value for specific method
	Value string
}

// Construct a new QueueSource, need no pointer
func NewQueueSource(queue string, config QueueSourceConfig) (*QueueSource, error) {
	conNew:=&(config)
	if conNew.Method==QUEUE_SOURCE_METHOD_API ||
		conNew.Method==QUEUE_SOURCE_METHOD_MYSQL ||
		conNew.Method==QUEUE_SOURCE_METHOD_FILE {

		conNew.Value=queue
	}else{
		return nil, errors.New("Unsupport QueueSource method.")
	}

	return &QueueSource{config:conNew}, nil
}

func NewQueueSourceByConfig(config *QueueSourceConfig) (*QueueSource, error) {
	if config.Method==QUEUE_SOURCE_METHOD_API ||
	config.Method==QUEUE_SOURCE_METHOD_MYSQL ||
	config.Method==QUEUE_SOURCE_METHOD_FILE {

		if strings.Trim(config.Value, " \t")=="" {
			return nil, errors.New("QueueSourceConfig Vaule field empty.")
		}
	}else{
		return nil, errors.New("Unsupport QueueSource method.")
	}

	return &QueueSource{config:config}, nil
}


//use cache first, update when needed
func (qs *QueueSource) GetDeviceQueue() {

}

func (qs *QueueSource) Cache() {

}

func (qs *QueueSource) Update() {

}





