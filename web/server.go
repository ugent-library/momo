package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/ugent-library/momo/engine"
	"github.com/ugent-library/momo/web/app"
	"github.com/ugent-library/momo/web/ui"
)

type Server interface {
	Start() error
}

type server struct {
	app  *app.App
	host string
	port int
}

type option func(*server)

func NewServer(e engine.Engine, opts ...option) Server {
	s := &server{
		app: app.New(e),
	}

	for _, opt := range opts {
		opt(s)
	}

	s.initRoutes()

	return s
}

func WithHost(h string) option {
	return func(s *server) {
		s.host = h
	}
}

func WithPort(p int) option {
	return func(s *server) {
		s.port = p
	}
}

func (s *server) initRoutes() {
	a := s.app
	r := a.Router

	// general middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// static file server
	r.Mount("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))

	// robots.txt
	r.Get("/robots.txt", Robots(a))

	for _, lens := range a.Lenses() {
		r.Route("/"+lens.Name, func(r chi.Router) {
			r.Use(app.ScopeSetter(lens.Scope))
			r.Use(app.ThemeSetter(lens.Theme))

			r.Get("/", ui.ListRecs(a))
			r.Get("/search", ui.SearchRecs(a))
			r.Get("/{id}", ui.ShowRec(a))
		})
	}
}

func (s *server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	fmt.Println(fmt.Sprintf("The momo server is running at http://%s.", addr))
	return http.ListenAndServe(addr, s.app.Router)
}
