package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/ugent-library/momo/http/ui"
	"github.com/ugent-library/momo/records"
)

type Lens struct {
	Name  string
	Scope records.Scope
	Theme string
}

type Options struct {
	Store       records.Storage
	SearchStore records.SearchStorage
	Host        string
	Port        int
}

func Start(opts Options) error {
	router := chi.NewRouter()

	router.Use(chimw.RequestID)
	router.Use(chimw.RealIP)
	router.Use(chimw.Logger)
	router.Use(chimw.Recoverer)

	for _, v := range loadLenses() {
		service := records.NewService(opts.Store, opts.SearchStore, v.Scope)
		handler := ui.NewSearch(service, v.Theme)
		router.Route("/v/"+v.Name, func(r chi.Router) {
			r.Get("/", handler.Index())
			r.Get("/search", handler.Search())
			r.Get("/{recID}", handler.Get())
		})
	}

	router.Mount("/s/", http.StripPrefix("/s/", http.FileServer(http.Dir("static"))))

	addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
	fmt.Println(fmt.Sprintf("The momo server is running at http://%s.", addr))
	return http.ListenAndServe(addr, router)
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
