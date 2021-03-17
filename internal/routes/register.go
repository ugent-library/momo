package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/ugent-library/momo/internal/controller"
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

	// robots.txt
	r.Get("/robots.txt", controller.Robots(e))

	for _, lens := range e.Lenses() {
		r.Route("/"+lens.Name, func(r chi.Router) {
			r.Use(mw.ScopeSetter(lens.Scope))
			r.Use(mw.ThemeSetter(lens.Theme))

			r.Get("/", recs.List)
			r.Get("/search", recs.Search)
			r.Get("/{id}", recs.Show)
		})
	}
}
