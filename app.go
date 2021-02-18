package momo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
)

type App struct {
	isDev         bool
	assetManifest map[string]string
	staticPath    string
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
	if a.isDev {
		a.loadAssetManifest()
	}
	path, ok := a.assetManifest[asset]
	if !ok {
		err := fmt.Errorf("Asset %s not found in manifest %s", asset, a.assetManifest)
		log.Println(err)
		return "", err
	}
	return path.Join(a.staticPath, path), nil
}

func (a *App) Start() {
	a.isDev = true
	a.staticPath = "/s/"
	a.loadAssetManifest()

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// rendering test
	funcs := template.FuncMap{
		"assetPath": a.assetPath,
	}
	tmpl := template.Must(template.New("layout.html").Funcs(funcs).ParseFiles("templates/layout.html"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, struct{ Title string }{"My Title"})
	})
	// rendering test

	r.Mount("/v/orpheus", ViewpointService{}.Handler())

	r.Mount(a.staticPath, http.StripPrefix(a.staticPath, http.FileServer(http.Dir("static"))))

	http.ListenAndServe(":3000", r)
}
