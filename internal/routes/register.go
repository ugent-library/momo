package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/ugent-library/momo/internal/controller"
	"github.com/ugent-library/momo/internal/ctx"
	"github.com/ugent-library/momo/internal/engine"
	mw "github.com/ugent-library/momo/internal/middleware"
)

func Register(r chi.Router, e engine.Engine) {
	recs := controller.NewRecs(e)

	// general middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// static file server
	r.Mount("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))

	// OAI-PMH provider
	r.Mount("/oai", controller.OAI(e))

	// robots.txt
	r.Get("/robots.txt", controller.Robots(e))

	uiRoutes := func(r chi.Router) {
		for _, collection := range e.Collections() {
			r.Route("/collection/"+collection.Name, func(r chi.Router) {
				r.Use(mw.SetLocale(e))
				r.Use(chimw.WithValue(ctx.CollectionKey, collection.Name))
				r.Use(chimw.WithValue(ctx.ThemeKey, collection.Theme))
				r.Get("/", recs.List)
				r.Get("/search", recs.Search)
				r.Get("/{id}", recs.Show)
			})
		}
	}

	for _, loc := range e.Locales() {
		r.Route("/"+loc.Language().String(), func(r chi.Router) {
			r.Use(chimw.WithValue(ctx.LocaleKey, loc))
			uiRoutes(r)
		})
	}

	uiRoutes(r)
}
