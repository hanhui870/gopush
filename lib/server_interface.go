package lib

const (
	HTTP_METHOD_GET = "GET"
	HTTP_METHOD_POST = "POST"
	HTTP_METHOD_PUT = "PUT"
)

type Server interface {
	GetTaskQueue() *TaskQueue

	GetEnv() EnvInfo
}