// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package api

import (
	"gopush/api/handler"
	"gopush/lib"
)

func NewApiV1Server(env lib.EnvInfo) *Server {
	server := NewServer(env)

	server.HandleFunc("/api/v1/send", handler.Send)
	server.HandleFunc("/api/v1/add-device", handler.AddDevice)

	return server
}