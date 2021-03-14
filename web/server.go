package web

import (
	"fmt"
	"net/http"

	"github.com/ugent-library/momo/web/app"
)

type Server struct {
	app  *app.App
	host string
	port int
}

type option func(*Server)

func NewServer(a *app.App, opts ...option) *Server {
	s := &Server{app: a}

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
	return http.ListenAndServe(addr, s.app.Router)
}
