package momo

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/ugent-library/momo/store"
)

type ViewpointHandler struct {
	funcs template.FuncMap
	store *store.Es
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
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(200)
		// TODO write template to a buffer first so we can show an error page
		if err := tmpl.ExecuteTemplate(w, "layout", data{Title: "Orpheus"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *ViewpointHandler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query().Get("q")
		hits, err := s.store.SearchRecs(q)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		renderJSON(w, 200, hits)
	}
}

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}
