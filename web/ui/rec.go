package ui

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/momo/engine"
	"github.com/ugent-library/momo/web/app"
)

func ListRecs(a *app.App, router chi.Router) {
	type data struct {
		Title string
	}
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		a.RenderHTML(w, r, http.StatusOK, "index", data{Title: "Search"})
	})
}

func GetRec(a *app.App, router chi.Router) {
	type data struct {
		Rec *engine.Rec
	}
	router.Get("/{recID}", func(w http.ResponseWriter, r *http.Request) {
		recID := chi.URLParam(r, "recID")
		rec, err := a.GetRec(recID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 404)
			return
		}
		a.RenderHTML(w, r, http.StatusOK, "show", data{Rec: rec})
	})
}

// TODO move route to api
func SearchRecs(a *app.App, router chi.Router) {
	router.Get("/search", func(w http.ResponseWriter, r *http.Request) {
		searchArgs := engine.SearchArgs{}
		err := a.DecodeForm(&searchArgs, r.URL.Query())

		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hits, err := a.SearchRecs(searchArgs.WithScope(app.GetScope(r)))
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		a.RenderJSON(w, r, http.StatusOK, hits)
	})
}
