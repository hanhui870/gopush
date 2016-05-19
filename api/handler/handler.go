// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package handler

import (
	"net/http"
	"fmt"
)

type PushApi struct {
	server Server
}

func NewPushApi(server Server) *PushApi {
	return &PushApi{server:server}
}

func (api *PushApi) Send(w http.ResponseWriter, r *http.Request) {
	formatNormalResponceHeader(w)

	fmt.Fprintln(w, `{"success":true,"message":"hello world."}`)
}

func (api *PushApi) AddDevice(w http.ResponseWriter, r *http.Request) {

}
