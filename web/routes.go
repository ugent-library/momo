package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/ugent-library/momo/web/app"
	"github.com/ugent-library/momo/web/ui"
)

func RegisterRoutes(a *app.App) {
	r := a.Router

	// general middleware
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// static file server
	r.Mount("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))

	// robots.txt
	r.Get("/robots.txt", Robots(a))

	for _, lens := range a.Lenses() {
		r.Route("/"+lens.Name, func(r chi.Router) {
			r.Use(app.ScopeSetter(lens.Scope))
			r.Use(app.ThemeSetter(lens.Theme))

			r.Get("/", ui.ListRecs(a))
			r.Get("/search", ui.SearchRecs(a))
			r.Get("/{id}", ui.ShowRec(a))
		})
	}
}
