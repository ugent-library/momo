package lens

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/form/v4"
	"github.com/ugent-library/momo/records"
	"github.com/unrolled/render"
)

type Handler struct {
	service     records.Service
	render      *render.Render
	formDecoder *form.Decoder
}

func NewHandler(service records.Service, layout string, funcs template.FuncMap) *Handler {
	if layout == "" {
		layout = "layout"
	}
	r := render.New(render.Options{
		Layout: layout,
		Funcs:  []template.FuncMap{funcs},
	})
	h := &Handler{
		service:     service,
		render:      r,
		formDecoder: form.NewDecoder(),
	}
	return h
}

func (s *Handler) Index() http.HandlerFunc {
	type data struct {
		Title string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.render.HTML(w, http.StatusOK, "index", data{Title: "Search"})
	}
}

func (s *Handler) Get() http.HandlerFunc {
	type data struct {
		Rec *records.Rec
	}
	return func(w http.ResponseWriter, r *http.Request) {
		recID := chi.URLParam(r, "recID")
		rec, err := s.service.GetRec(recID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 404)
			return
		}
		s.render.HTML(w, http.StatusOK, "show", data{Rec: rec})
	}
}

// TODO move route to api
func (s *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchArgs := records.SearchArgs{}
		err := s.formDecoder.Decode(&searchArgs, r.URL.Query())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hits, err := s.service.SearchRecs(searchArgs)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.render.JSON(w, http.StatusOK, hits)
	}
}
