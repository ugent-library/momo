package ui

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path"

	"github.com/gorilla/schema"
	"github.com/ugent-library/momo/listing"
)

type ViewpointHandler struct {
	listingService listing.Service
	layout         string
	funcs          template.FuncMap
}

func (s *ViewpointHandler) Index() http.HandlerFunc {
	layout := s.layout
	if layout == "" {
		layout = "layout.tmpl"
	}
	tmpl, err := template.New("layout.tmpl").Funcs(s.funcs).ParseFiles(path.Join("templates", layout), "templates/index.tmpl")
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
		if err := tmpl.ExecuteTemplate(w, "layout", data{Title: "Search"}); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (s *ViewpointHandler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchArgs := listing.SearchArgs{}
		decoder := schema.NewDecoder()
		err := decoder.Decode(&searchArgs, r.URL.Query())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hits, err := s.listingService.SearchRecs(searchArgs)
		if err != nil {
			log.Println(err)
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
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}
