package ui

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/momo/engine"
	"github.com/ugent-library/momo/web/app"
)

func ListRecs(a *app.App) http.HandlerFunc {
	type data struct {
		Title string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		a.RenderHTML(w, r, "index", data{Title: "Search"})
	}
}

func ShowRec(a *app.App) http.HandlerFunc {
	type data struct {
		Rec *engine.Rec
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		rec, err := a.GetRec(id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 404)
			return
		}

		a.RenderHTML(w, r, "show", data{Rec: rec})
	}
}

// TODO move to api
func SearchRecs(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		a.RenderJSON(w, r, hits)
	}
}
