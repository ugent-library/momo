package momo

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/unrolled/render"
)

type ViewpointHandler struct {
	r    *render.Render
	recs *Recs
}

func (s *ViewpointHandler) Handler() http.Handler {

	r := chi.NewRouter()
	r.Get("/", s.Index())
	r.Get("/search", s.Search())
	return r
}

func (s *ViewpointHandler) Index() http.HandlerFunc {
	type data struct {
		Title string
		Hits  map[string]interface{}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		hits := s.recs.AutocompleteSearch(q)
		s.r.HTML(w, http.StatusOK, "orpheus/index", data{Title: "Orpheus", Hits: hits})
	}
}

func (s *ViewpointHandler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		hits := s.recs.AutocompleteSearch(q)
		s.r.JSON(w, http.StatusOK, hits)
	}
}
