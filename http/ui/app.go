package ui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/ugent-library/momo/http/ui/lens"
	"github.com/ugent-library/momo/records"
)

type Lens struct {
	Name   string
	Scope  records.Scope
	Layout string
}

type App struct {
	store         records.Storage
	searchStore   records.SearchStorage
	Port          int
	assetManifest map[string]string
	staticPath    string
	funcs         template.FuncMap
}

func New(store records.Storage, searchStore records.SearchStorage) *App {
	a := &App{
		store:       store,
		searchStore: searchStore,
		staticPath:  "/s/",
	}
	a.loadAssetManifest()
	a.funcs = template.FuncMap{
		"assetPath": a.assetPath,
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
		handler := lens.NewHandler(service, v.Layout, a.funcs)
		r.Route("/v/"+v.Name, func(r chi.Router) {
			r.Get("/", handler.Index())
			r.Get("/search", handler.Search())
			r.Get("/{recID}", handler.Get())
		})
	}

	r.Mount(a.staticPath, http.StripPrefix(a.staticPath, http.FileServer(http.Dir("static"))))

	return r
}

func (a *App) loadAssetManifest() {
	path := "static/mix-manifest.json"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Couldn't read %s: %s", path, err)
	}
	manifest := make(map[string]string)
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatalf("Couldn't parse %s: %s", path, err)
	}
	a.assetManifest = manifest
}

func (a *App) assetPath(asset string) (string, error) {
	p, ok := a.assetManifest[asset]
	if !ok {
		err := fmt.Errorf("Asset %s not found in manifest %s", asset, a.assetManifest)
		log.Println(err)
		return "", err
	}
	return path.Join(a.staticPath, p), nil
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
