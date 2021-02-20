package momo

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/go-chi/chi"
	chimw "github.com/go-chi/chi/middleware"
	"github.com/unrolled/render"
)

type App struct {
	isDev         bool
	assetManifest map[string]string
	staticPath    string
	r             *render.Render
}

// TODO mutex
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
	p, ok := a.assetManifest[asset]
	if !ok {
		err := fmt.Errorf("Asset %s not found in manifest %s", asset, a.assetManifest)
		log.Println(err)
		return "", err
	}
	return path.Join(a.staticPath, p), nil
}

func (a *App) Start() {
	a.isDev = true
	a.staticPath = "/s/"
	a.r = render.New(render.Options{
		IsDevelopment: a.isDev,
		Layout:        "layout",
		Funcs: []template.FuncMap{template.FuncMap{
			"assetPath": a.assetPath,
		}},
		XMLContentType: "application/xml",
	})
	a.loadAssetManifest()

	es, err := elasticsearch.NewDefaultClient()
	if err == nil {
		// log.Println(elasticsearch.Version)
		// log.Println(es.Info())
	} else {
		log.Fatalf("Can't create es client: %s", err)
	}

	recs := &Recs{
		es:    es,
		index: "momo_rec",
	}

	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// rendering test
	// funcs := template.FuncMap{
	// 	"assetPath": a.assetPath,
	// }
	// tmpl := template.Must(template.New("layout.html").Funcs(funcs).ParseFiles("templates/layout.html"))
	// r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl.Execute(w, struct{ Title string }{"My Title"})
	// })
	// rendering test

	r.Mount("/v/orpheus", (&ViewpointHandler{r: a.r, recs: recs}).Handler())

	r.Mount(a.staticPath, http.StripPrefix(a.staticPath, http.FileServer(http.Dir("static"))))

	fmt.Println("The momo server is running at http://localhost:3000.")
	http.ListenAndServe("localhost:3000", r)
}
