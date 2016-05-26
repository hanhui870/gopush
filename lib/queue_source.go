package lib

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

	ApiUri   string
}

//use cache first, update when needed
func (qs *QueueSource) GetDeviceQueue() {

}

func (qs *QueueSource) Cache() {

}

func (qs *QueueSource) Update() {

}





