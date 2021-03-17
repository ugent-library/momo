package render

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/ugent-library/momo/web/ctx"
	"github.com/ugent-library/momo/web/theme"
)

type Data struct {
	Yield interface{}
}

type View struct {
	layout    string
	templates map[string]*template.Template
}

func NewView(layout string, files []string, funcs ...template.FuncMap) *View {
	templates := make(map[string]*template.Template)

	for _, t := range theme.Themes() {
		templateFiles := make([]string, len(files))
		for i, f := range files {
			templateFiles[i] = "themes/" + t.Name() + "/templates/" + f + ".gohtml"
		}
		layoutFiles, err := filepath.Glob("themes/" + t.Name() + "/templates/layouts/*.gohtml")
		if err != nil {
			panic(err)
		}
		templateFiles = append(templateFiles, layoutFiles...)

		tmpl := template.New("")
		if f := t.Funcs(); f != nil {
			tmpl.Funcs(f)
		}
		for _, f := range funcs {
			tmpl.Funcs(f)
		}
		tmpl, err = tmpl.ParseFiles(templateFiles...)
		if err != nil {
			panic(err)
		}

		templates[t.Name()] = tmpl
	}

	return &View{
		templates: templates,
		layout:    layout,
	}
}

func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is used to render the view with the predefined layout.
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(ctx.GetStatus(r))

	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
		// do nothing
	default:
		vd = Data{
			Yield: data,
		}
	}

	themeName := ctx.GetTheme(r)
	tmpl, found := v.templates[themeName]
	if !found {
		err := errors.New("template not found for theme " + themeName)
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, v.layout, vd); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}
