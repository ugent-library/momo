package view

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	"html/template"
	"github.com/unrolled/render"
)

type renderer struct {
	Options       render.Options
	Theme         string
	AssetManifest map[string]string
	StaticPath    string
	Funcs         template.FuncMap
}

type Renderer interface {
	Create() (*render.Render)
}

func NewRenderer(theme string) Renderer {
	r := &renderer{}
	if theme == "" {
		theme = "opale" // default theme: Opale
	}
	r.Funcs = template.FuncMap{
		"assetPath": r.assetPath,
	}
	options := render.Options{
		Directory: fmt.Sprintf("themes/%s/templates", theme),
		Layout: "layout",
		Funcs:  []template.FuncMap{r.Funcs},
	}
	r.Options = options
	r.Theme = theme
	r.StaticPath = "/s/"
	r.loadAssetManifest()
	return r
}

func (r *renderer) Create() (*render.Render) {
	return render.New(r.Options)
}

func (r *renderer) loadAssetManifest() {
	path := fmt.Sprintf("static/%s/mix-manifest.json", r.Theme)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Couldn't read %s: %s", path, err)
	}
	manifest := make(map[string]string)
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		log.Fatalf("Couldn't parse %s: %s", path, err)
	}
	r.AssetManifest = manifest
}

func (r *renderer) assetPath(asset string) (string, error) {
	p, ok := r.AssetManifest[asset]
	if !ok {
		err := fmt.Errorf("Asset %s not found in manifest %s", asset, r.AssetManifest)
		log.Println(err)
		return "", err
	}
	return path.Join(r.StaticPath, p), nil
}
