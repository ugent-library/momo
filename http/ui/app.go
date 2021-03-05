package ui

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/ugent-library/momo/http/ui/lens"
	"github.com/ugent-library/momo/records"
)

type Lens struct {
	Name   string
	Scope  records.Scope
	Layout string
	Theme  string
}

type App struct {
	store         records.Storage
	searchStore   records.SearchStorage
	Port          int
	staticPath    string
}

func New(store records.Storage, searchStore records.SearchStorage) *App {
	a := &App{
		store:       store,
		searchStore: searchStore,
		staticPath:  "/s/",
	}
	return a
}
func (a *App) Start() {
	fmt.Println(fmt.Sprintf("The momo server is running at http://localhost:%d.", a.Port))
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", a.Port), a.router())
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) router() *chi.Mux {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	for _, v := range loadLenses() {
		service := records.NewService(a.store, a.searchStore, v.Scope)
		handler := lens.NewHandler(service, v.Theme)
		r.Route("/v/"+v.Name, func(r chi.Router) {
			r.Get("/", handler.Index())
			r.Get("/search", handler.Search())
			r.Get("/{recID}", handler.Get())
		})
	}

	r.Mount(a.staticPath, http.StripPrefix(a.staticPath, http.FileServer(http.Dir("static"))))

	return r
}

func loadLenses() []Lens {
	jsonFile, err := os.Open("etc/lenses.json")
	defer jsonFile.Close()
	if err != nil {
		log.Fatal(err)
	}
	v := make([]Lens, 0)
	if err := json.NewDecoder(jsonFile).Decode(&v); err != nil {
		log.Fatal(err)
	}
	return v
}
