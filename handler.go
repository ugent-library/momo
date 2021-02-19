package momo

import (
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
	return r
}

func (s *ViewpointHandler) Index() http.HandlerFunc {
	tmpl, err := template.New("index.html").Funcs(s.funcs).ParseFiles("templates/orpheus/index.html")
	if err != nil {
		log.Printf("error parsing template: %s", err)
	}
	type data struct {
		Title string
		Hits  map[string]interface{}
	}
	return func(w http.ResponseWriter, r *http.Request) {
		qs := r.URL.Query().Get("qs")
		err := tmpl.Execute(w, data{Title: "Orpheus", Hits: s.recs.AutocompleteSearch(qs)})
		if err != nil {
			log.Printf("error executing template: %s", err)
		}
	}
}
