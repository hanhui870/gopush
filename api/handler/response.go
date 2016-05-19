package handler

import ()

type Response struct {
	//fail or success
	Error   bool `json:"error"`
	//system message
	Message string `json:"message"`
	//error code, 0 of success
	Code    int `json:"code"`
}

type SendResponse struct {
	Response

	PushID int64 `json:"push-id"`
}


