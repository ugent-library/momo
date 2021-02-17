package momo

import (
	"net/http"
	"text/template"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/go-webpack/webpack"
)

type App struct {
}

func (a App) Start() {
	isDev := true
	webpack.FsPath = "static"
	webpack.Init(isDev)

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	funcs := template.FuncMap{"asset": webpack.AssetHelper}
	tmpl := template.Must(template.New("layout.html").Funcs(funcs).ParseFiles("templates/layout.html"))
	// tmpl := template.Must(template.ParseFiles("templates/layout.html"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, struct{ Title string }{"My Title"})
	})

	r.Mount("/v/orpheus", ViewpointService{}.Handler())

	r.Mount("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":3000", r)
}
