package handler

import (
	"gopush/lib"
)

const (
	HTTP_METHOD_GET = "GET"
	HTTP_METHOD_POST = "POST"
	HTTP_METHOD_PUT = "PUT"

	API_CODE_OK = iota

//post method
	API_CODE_POST_NEEDED
//param required
	API_CODE_PARAM_REQUIRED
	API_CODE_PARAM_ERROR
	API_CODE_QUEUE_BUILD
	API_CODE_TASK_ERROR

	DEVICEID_SEP = ","
)

type Server interface {
	GetTaskQueue() *lib.TaskQueue

	GetEnv() lib.EnvInfo
}