// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package api

import (
	"net/http"

	"gopush/lib"
)

//An http server type
type Server struct {
	handler *http.ServeMux

	server  *http.Server

	//Server task queue
	task    *lib.TaskQueue

	//Env info
	env     lib.EnvInfo
}

func NewServer(env lib.EnvInfo) *Server {
	handle := http.NewServeMux()
	return &Server{handler:handle, server:&http.Server{Handler:handle}, env:env, task:lib.NewTaskQueue()}
}

func (s *Server) Start() error {
	s.server.Addr = s.env.GetServerAddr()

	s.env.GetLogger().Println("Server http://" + s.server.Addr + " started...")

	//Server taskqueue run
	go s.task.Run()

	return s.server.ListenAndServe()
}

func (s *Server) Stop() {

}

func (s *Server) GetTaskQueue() *lib.TaskQueue {
	return s.task
}

func (s *Server) GetEnv() lib.EnvInfo {
	return s.env
}

// Handle registers th+e handler for the given pattern in the DefaultServeMux.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.handler.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern in the DefaultServeMux.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.handler.HandleFunc(pattern, handler)
}

