package momo

import (
	"net/http"

	"github.com/go-chi/chi"
)

type ViewpointService struct{}

func (s ViewpointService) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", s.Index)
	return r
}

func (s ViewpointService) Index(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("viewpoint"))
}
