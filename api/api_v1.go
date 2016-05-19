// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package api

import (
	"gopush/api/handler"
	"gopush/lib"
)

func NewApiV1Server(env lib.EnvInfo) *Server {
	server := NewServer(env)

	api := handler.NewPushApi(server)
	server.HandleFunc("/api/v1/send", api.Send)
	server.HandleFunc("/api/v1/add-device", api.AddDevice)

	return server
}