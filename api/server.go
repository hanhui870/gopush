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

	//Worker pool
	pool    *lib.Pool

	//Env info
	env     lib.EnvInfo
}

func NewServer(env lib.EnvInfo) *Server {
	handle := http.NewServeMux()
	return &Server{handler:handle, server:&http.Server{Handler:handle}, env:env}
}

func (s *Server) Start() error {
	s.server.Addr = s.env.GetServerAddr()

	s.env.GetLogger().Println("Server http://" + s.server.Addr + " started...")

	//pool worker run
	go s.pool.Run()

	return s.server.ListenAndServe()
}

func (s *Server) Stop() {

}

func (s *Server) GetPool() *lib.Pool {
	return s.pool
}

func (s *Server) GetEnv() lib.EnvInfo {
	return s.env
}

func (s *Server) SetPool(pool *lib.Pool) (bool) {
	s.pool = pool

	return true
}

// Handle registers th+e handler for the given pattern in the DefaultServeMux.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.handler.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern in the DefaultServeMux.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.handler.HandleFunc(pattern, handler)
}

