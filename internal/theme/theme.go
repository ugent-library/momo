package theme

import (
	"fmt"
	"html/template"

	"github.com/ugent-library/go-mix/mix"
)

var (
	themes = make(map[string]Theme)
)

type Theme interface {
	Name() string
	Funcs() template.FuncMap
}

type theme struct {
	name  string
	funcs template.FuncMap
}

func Register(t Theme) {
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
	list := make([]Theme, 0, len(themes))
	for _, theme := range themes {
		list = append(list, theme)
	}
	return list
}

func New(name string) Theme {
	t := &theme{
		name: name,
	}
	t.funcs = mix.FuncMap(mix.Config{
		ManifestFile: fmt.Sprintf("static/%s/mix-manifest.json", name),
		PublicPath:   fmt.Sprintf("/s/%s/", name),
	})

	return t
}

func (t *theme) Name() string {
	return t.name
}

func (t *theme) Funcs() template.FuncMap {
	return t.funcs
}
