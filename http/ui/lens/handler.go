package lens

import (
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/form/v4"
	"github.com/ugent-library/momo/listing"
	"github.com/unrolled/render"
)

type Handler struct {
	listingService listing.Service
	render         *render.Render
	formDecoder    *form.Decoder
}

func NewHandler(listingService listing.Service, layout string, funcs template.FuncMap) *Handler {
	if layout == "" {
		layout = "layout"
	}
	r := render.New(render.Options{
		Layout: layout,
		Funcs:  []template.FuncMap{funcs},
	})
	h := &Handler{
		listingService: listingService,
		render:         r,
		formDecoder:    form.NewDecoder(),
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
		searchArgs := listing.SearchArgs{}
		err := s.formDecoder.Decode(&searchArgs, r.URL.Query())
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
		s.render.JSON(w, http.StatusOK, hits)
	}
}
