// Copyright 2016 祝景法(Bruce)@haimi.com. www.haimi.com All rights reserved.
package api

import (
	"net/http"

	"gopush/lib"
)

//An http server type
type Server struct {
	server *http.ServeMux

	//Worker pool
	pool   lib.Pool

	//Env info
	env    lib.EnvInfo
}

func NewServer() *Server {
	return &Server{server:http.NewServeMux()}
}

func (s *Server) Start() {

}

func (s *Server) Stop() {

}

// Handle registers the handler for the given pattern in the DefaultServeMux.
func (s *Server) Handle(pattern string, handler http.Handler) {
	s.server.Handle(pattern, handler)
}

// HandleFunc registers the handler function for the given pattern in the DefaultServeMux.
func (s *Server) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	s.server.HandleFunc(pattern, handler)
}

