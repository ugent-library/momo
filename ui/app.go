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

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/ugent-library/momo/listing"
	"github.com/ugent-library/momo/storage"
	"github.com/ugent-library/momo/ui/lens"
)

type Lens struct {
	Name        string
	SearchScope listing.SearchScope
	Layout      string
}

type App struct {
	assetManifest map[string]string
	staticPath    string
	funcs         template.FuncMap
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

func (a *App) Start() {
	a.staticPath = "/s/"
	a.funcs = template.FuncMap{
		"assetPath": a.assetPath,
	}
	a.loadAssetManifest()

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatal(err)
	}
	mapping, err := ioutil.ReadFile("etc/es6/rec_mapping.json")
	if err != nil {
		log.Fatal(err)
	}
	store := &storage.Es6{
		Client:       client,
		IndexName:    "momo_rec",
		IndexMapping: string(mapping),
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	for _, v := range loadLenses() {
		listingService := listing.NewService(store, v.SearchScope)
		handler := lens.NewHandler(listingService, v.Layout, a.funcs)
		r.Route("/v/"+v.Name, func(r chi.Router) {
			r.Get("/", handler.Index())
			r.Get("/search", handler.Search())
		})
	}

	r.Mount(a.staticPath, http.StripPrefix(a.staticPath, http.FileServer(http.Dir("static"))))

	fmt.Println("The momo server is running at http://localhost:3000.")
	http.ListenAndServe("localhost:3000", r)
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
