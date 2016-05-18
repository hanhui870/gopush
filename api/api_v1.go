// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package api

import (
	"gopush/api/handler"
)

func NewApiV1Server() *Server {
	server := NewServer()

	server.HandleFunc("/api/v1/send", handler.Send)
	server.HandleFunc("/api/v1/add-device", handler.AddDevice)

	return server
}