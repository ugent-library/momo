package controller

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/momo/internal/ctx"
	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/form"
	"github.com/ugent-library/momo/internal/metadata"
	"github.com/ugent-library/momo/internal/render"
)

type RecController struct {
	engine   engine.Engine
	listView *render.View
	showView *render.View
}

func NewRecController(e engine.Engine) *RecController {
	return &RecController{
		engine:   e,
		listView: render.NewView(e, "app", []string{"rec/list"}),
		showView: render.NewView(e, "app", []string{"rec/show"}, template.FuncMap{
			"renderSourceView":   renderSourceView,
			"renderInternalView": renderInternalView,
			"renderRepresentation": func(rec *engine.Rec, format string) template.HTML {
				rep, err := e.GetRepresentation(rec.ID, format)
				if err != nil {
					return template.HTML("")
				}
				return template.HTML(string(rep.Data))
			},
		}),
	}
}

func (c *RecController) List(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Title string
	}
	c.listView.Render(w, r, data{Title: "Search"})
}

func (c *RecController) Show(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Rec metadata.Rec
	}

	id := chi.URLParam(r, "id")
	rec, err := c.engine.GetRec(id)
	if err != nil || rec.Collection != ctx.GetCollection(r) {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}

	c.showView.Render(w, r, data{Rec: metadata.WrapRec(rec)})
}

func (c *RecController) Search(w http.ResponseWriter, r *http.Request) {
	args := engine.SearchArgs{Facets: []string{"type"}, Highlight: true}
	err := form.Decode(&args, r.URL.Query())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if args.Size == 0 {
		args.Size = 10
	}
	hits, err := c.engine.SearchRecs(args.WithFilter("collection", ctx.GetCollection(r)))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, hits)
}

func renderSourceView(rec *engine.Rec) template.HTML {
	var b strings.Builder

	if rec.SourceFormat == "marcinjson" && rec.SourceMetadata != nil {
		marc := gjson.ParseBytes([]byte(rec.SourceMetadata))

		b.WriteString(`<table class="table table-sm table-striped">`)
		b.WriteString(`<tr><th colspan="4">` + marc.Get(`leader`).String() + `</th></tr>`)
		marc.Get(`fields`).ForEach(func(_, field gjson.Result) bool {
			field.ForEach(func(code, f gjson.Result) bool {
				b.WriteString(`<tr><th class="table-active">` + code.String() + `</th>`)
				if f.IsObject() {
					b.WriteString(`<td>` + f.Get(`ind1`).String() + `</td>`)
					b.WriteString(`<td>` + f.Get(`ind2`).String() + `</td>`)
					b.WriteString(`<td>`)
					f.Get(`subfields`).ForEach(func(_, subfield gjson.Result) bool {
						subfield.ForEach(func(code, sf gjson.Result) bool {
							b.WriteString(`<span class="text-muted">` + code.String() + `</span> ` + sf.String() + ` `)
							return false
						})
						return true
					})
					b.WriteString(`</td>`)
				} else {
					b.WriteString(`<td colspan="3">` + f.String() + `</td>`)
				}
				b.WriteString(`</tr>`)
				return false
			})
			return true
		})

		b.WriteString(`</table>`)
	}

	return template.HTML(b.String())
}

func renderInternalView(rec *engine.Rec) template.HTML {
	b, _ := json.MarshalIndent(rec.Metadata, "", "\t")
	return template.HTML("<code><pre>" + string(b) + "</pre></code>")
}
