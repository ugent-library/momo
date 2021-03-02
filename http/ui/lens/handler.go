package lens

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/ugent-library/momo/records"
	"github.com/unrolled/render"
)

type Handler struct {
	searchService records.SearchService
	render        *render.Render
	formDecoder   *form.Decoder
}

func NewHandler(searchService records.SearchService, layout string, funcs template.FuncMap) *Handler {
	if layout == "" {
		layout = "layout"
	}
	r := render.New(render.Options{
		Layout: layout,
		Funcs:  []template.FuncMap{funcs},
	})
	h := &Handler{
		searchService: searchService,
		render:        r,
		formDecoder:   form.NewDecoder(),
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

func (s *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchArgs := records.SearchArgs{}
		err := s.formDecoder.Decode(&searchArgs, r.URL.Query())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hits, err := s.searchService.Search(searchArgs)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.render.JSON(w, http.StatusOK, hits)
	}
}
