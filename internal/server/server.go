package server

import (
	"fmt"
	"net/http"
)

type Server struct {
	handler http.Handler
	host    string
	port    int
}

type option func(*Server)

func New(h http.Handler, opts ...option) *Server {
	s := &Server{handler: h}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

func WithHost(h string) option {
	return func(s *Server) {
		s.host = h
	}
}

func WithPort(p int) option {
	return func(s *Server) {
		s.port = p
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	fmt.Println(fmt.Sprintf("The momo server is running at http://%s.", addr))
	return http.ListenAndServe(addr, s.handler)
}
