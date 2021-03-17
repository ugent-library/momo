package controller

import (
	"bytes"
	"encoding/json"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/momo/internal/ctx"
	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/form"
	"github.com/ugent-library/momo/internal/render"
)

type Recs struct {
	engine   engine.Engine
	listView *render.View
	showView *render.View
}

func NewRecs(e engine.Engine) *Recs {
	return &Recs{
		engine:   e,
		listView: render.NewView("app", []string{"rec/list"}),
		showView: render.NewView("app", []string{"rec/show"}, template.FuncMap{
			"renderTitle":        renderTitle,
			"renderMetadata":     renderMetadata,
			"renderSourceView":   renderSourceView,
			"renderInternalView": renderInternalView,
		}),
	}
}

func (c *Recs) List(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Title string
	}
	c.listView.Render(w, r, data{Title: "Search"})
}

func (c *Recs) Show(w http.ResponseWriter, r *http.Request) {
	type data struct {
		Rec *engine.Rec
	}

	id := chi.URLParam(r, "id")
	rec, err := c.engine.GetRec(id)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), 404)
		return
	}

	c.showView.Render(w, r, data{Rec: rec})
}

func (c *Recs) Search(w http.ResponseWriter, r *http.Request) {
	searchArgs := engine.SearchArgs{}
	err := form.Decode(&searchArgs, r.URL.Query())

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hits, err := c.engine.SearchRecs(searchArgs.WithScope(ctx.GetScope(r)))
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, hits)
}

func renderTitle(j json.RawMessage) template.HTML {
	title := gjson.GetBytes(j, "title").String()
	return template.HTML(title)
}

func renderMetadata(j json.RawMessage) template.HTML {
	var b strings.Builder
	m := gjson.ParseBytes(j)
	b.WriteString("<dl>")
	renderMetadataField(&b, &m, "author.#.name", "Authors")
	renderMetadataField(&b, &m, "abstract.#.text", "Abstract")
	renderMetadataField(&b, &m, "edition", "Edition")
	renderMetadataField(&b, &m, "publisher", "Publisher")
	renderMetadataField(&b, &m, "placeOfPublication", "Place of publication")
	renderMetadataField(&b, &m, "publicationDate", "Date published")
	renderMetadataField(&b, &m, "doi", "DOI")
	renderMetadataField(&b, &m, "isbn", "ISBN")
	renderMetadataField(&b, &m, "note.#.text", "Note")
	b.WriteString("</dl>")
	return template.HTML(b.String())
}

func renderMetadataField(b *strings.Builder, m *gjson.Result, path string, dt string) {
	res := m.Get(path)
	if res.Exists() {
		b.WriteString("<dt>" + html.EscapeString(dt) + "</dt>")
		res.ForEach(func(_, v gjson.Result) bool {
			b.WriteString("<dd>" + html.EscapeString(v.String()) + "</dd>")
			return true
		})
	}
}

func renderSourceView(j json.RawMessage) template.HTML {
	var b strings.Builder
	m := gjson.ParseBytes(j)
	marc := m.Get(`@this.#(metadata_format=="marc-in-json").metadata`)
	if marc.Exists() {
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
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, rec.RawMetadata, "", "\t")
	return template.HTML("<code><pre>" + prettyJSON.String() + "</pre></code>")
}
