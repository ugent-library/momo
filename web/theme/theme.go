package theme

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"path"
	"sync"
)

var (
	themesMu sync.RWMutex
	themes   = make(map[string]Theme)
)

type Theme interface {
	Name() string
	FuncMaps() []template.FuncMap
}

type theme struct {
	name          string
	funcMaps      []template.FuncMap
	assetManifest map[string]string
}

func Register(t Theme) {
	themesMu.Lock()
	defer themesMu.Unlock()
	if t == nil {
		panic("momo: theme is nil")
	}
	if t.Name() == "" {
		panic("momo: theme name is empty")
	}
	if _, dup := themes[t.Name()]; dup {
		panic("momo: theme " + t.Name() + " already exists")
	}
	themes[t.Name()] = t
}

// Themes returns a list of registered themes.
func Themes() []Theme {
	themesMu.Lock()
	defer themesMu.Unlock()
	list := make([]Theme, 0, len(themes))
	for _, theme := range themes {
		list = append(list, theme)
	}
	return list
}

func New(name string) Theme {
	t := &theme{
		name:          name,
		assetManifest: loadAssetManifest(name),
	}
	t.funcMaps = []template.FuncMap{{
		"assetPath": t.assetPath,
	}}
	return t
}

func (t *theme) Name() string {
	return t.name
}

func (t *theme) FuncMaps() []template.FuncMap {
	return t.funcMaps
}

func (t *theme) assetPath(asset string) (string, error) {
	p, ok := t.assetManifest[asset]
	if !ok {
		err := fmt.Errorf("Asset %s not found in manifest %s", asset, t.assetManifest)
		log.Println(err)
		return "", err
	}
	return path.Join("/s/", t.name, p), nil
}

func loadAssetManifest(name string) (manifest map[string]string) {
	path := fmt.Sprintf("static/%s/mix-manifest.json", name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Couldn't read %s: %s", path, err)
	}
	manifest = make(map[string]string)
	if err = json.Unmarshal(data, &manifest); err != nil {
		log.Fatalf("Couldn't parse %s: %s", path, err)
	}
	return
}
