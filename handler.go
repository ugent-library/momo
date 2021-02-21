package momo

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type ViewpointHandler struct {
	funcs template.FuncMap
	recs  *Recs
}

func (s *ViewpointHandler) Handler() http.Handler {

	r := chi.NewRouter()
	r.Get("/", s.Index())
	r.Get("/search", s.Search())
	return r
}

func (s *ViewpointHandler) Index() http.HandlerFunc {
	tmpl, err := template.New("layout.tmpl").Funcs(s.funcs).ParseFiles("templates/layout.tmpl", "templates/orpheus/index.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	type data struct {
		Title string
		Hits  map[string]interface{}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		hits := s.recs.AutocompleteSearch(q)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		if err := tmpl.ExecuteTemplate(w, "layout", data{Title: "Orpheus", Hits: hits}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *ViewpointHandler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		hits := s.recs.AutocompleteSearch(q)

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(true)
		if err := enc.Encode(hits); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write(buf.Bytes())
	}
}
