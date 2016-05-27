// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package handler

import (
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
	"strings"
	"gopush/lib"

	"github.com/twinj/uuid"
	"strconv"
)

type PushApi struct {
	server lib.Server
}

func NewPushApi(server lib.Server) *PushApi {
	return &PushApi{server:server}
}

// Send API
//
// DESC: Send an notification to the pool
// Params:
//		title: notification title
//		body: notification body info
//		custom: json string, map[string][string], eg. custom={"payload": "haimi-590"}
//		sound: notification sound
//		queue: send queue, empty will use default all users.
//			depends on runtime/config/config.ini queue.method value, file, sql, api has different meanings.
//		deviceids: Send to specified id, not required. delimited by ","
func (api *PushApi) Send(w http.ResponseWriter, r *http.Request) {
	formatNormalResponceHeader(w)
	r.ParseForm()

	api.server.GetEnv().GetLogger().Println("Receive request: ", r.Form)

	if r.Method != lib.HTTP_METHOD_POST {
		api.OutputResponse(w, &Response{Error:true, Message:"HTTP method POST is required.", Code:API_CODE_POST_NEEDED})
		return
	}

	title, err := GetParamString(r, "title")
	if err != nil {
		api.OutputResponse(w, &Response{Error:true, Message:"Param title is required.", Code:API_CODE_PARAM_REQUIRED})
		return
	}

	body, err := GetParamString(r, "body")
	if err != nil {
		api.OutputResponse(w, &Response{Error:true, Message:"Param body is required.", Code:API_CODE_PARAM_REQUIRED})
		return
	}

	tmpArr, err := GetParamString(r, "custom")
	var custom map[string]string
	custom = make(map[string]string, 100)
	if err == nil {
		err = json.Unmarshal(bytes.NewBufferString(tmpArr).Bytes(), &custom)
		if err != nil {
			api.OutputResponse(w, &Response{Error:true, Message:"Param custom json parse failed:" + err.Error(), Code:API_CODE_PARAM_ERROR})
			return
		}
	}

	str, err := GetParamString(r, "sound")
	var sound string
	if err == nil {
		sound = str
	}else {
		sound = ""
	}

	str, err = GetParamString(r, "queue")
	var queue string
	if err == nil {
		queue = str
	}else {
		queue = ""
	}

	str, err = GetParamString(r, "deviceids")
	var deviceids []string
	if err == nil {
		deviceids = strings.Split(str, DEVICEID_SEP)
	}else {
		deviceids = nil
	}

	msg := &lib.Message{Title:title, Body:body, Sound:sound, Custom:custom, Uuid:uuid.NewV1().String()}

	qb := lib.NewQueueBuilder(queue, deviceids)
	devicequeue, err := qb.ToDeviceQueue(api.server.GetEnv().GetPoolConfig().Capacity)
	if err != nil {
		api.OutputResponse(w, &Response{Error:true, Message:"Build send queue failed:" + err.Error(), Code:API_CODE_QUEUE_BUILD})
		return
	}

	position, err := api.server.GetTaskQueue().Add(devicequeue, msg)
	if err != nil {
		api.OutputResponse(w, &Response{Error:true, Message:"Add to taskqueue error:" + err.Error(), Code:API_CODE_TASK_ERROR})
		return
	}

	//Send Response
	resp := new(SendResponse)
	resp.Position = position
	resp.PushID = msg.Uuid
	resp.Error = false
	resp.Message = "Sent:" + msg.Uuid + " Position:" + strconv.Itoa(position)
	resp.Code = API_CODE_OK
	api.OutputResponse(w, resp)
	return
}

func (api *PushApi) AddDevice(w http.ResponseWriter, r *http.Request) {

}

func (api *PushApi) FormatResponseJson(resp interface{}) (string, error) {
	result, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}else {
		return string(result), nil
	}
}

func (api *PushApi) OutputResponse(w http.ResponseWriter, resp interface{}) {
	resp, err := api.FormatResponseJson(resp)
	if err == nil {
		fmt.Fprintln(w, resp)
		api.server.GetEnv().GetLogger().Println("Resp:", resp)
	}else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
		api.server.GetEnv().GetLogger().Println("Found error:", err)
	}
}
