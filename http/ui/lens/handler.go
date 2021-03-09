package lens

import (
	"encoding/json"
	"html"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-playground/form/v4"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/momo/records"
	"github.com/ugent-library/momo/view"
	"github.com/unrolled/render"
)

type Handler struct {
	service     records.Service
	render      *render.Render
	formDecoder *form.Decoder
}

func NewHandler(service records.Service, theme string) *Handler {
	funcs := template.FuncMap{
		"renderMetadata": renderMetadata,
		"renderSource":   renderSource,
	}
	r := view.NewRenderer(theme, funcs)
	h := &Handler{
		service:     service,
		render:      r.Create(),
		formDecoder: form.NewDecoder(),
	}
	return h
}

func (s *Handler) Index() http.HandlerFunc {
	type data struct {
		Title string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		s.render.HTML(w, http.StatusOK, "index", data{Title: "Search"})
	}
}

func (s *Handler) Get() http.HandlerFunc {
	type data struct {
		Rec *records.Rec
	}
	return func(w http.ResponseWriter, r *http.Request) {
		recID := chi.URLParam(r, "recID")
		rec, err := s.service.GetRec(recID)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 404)
			return
		}
		s.render.HTML(w, http.StatusOK, "show", data{Rec: rec})
	}
}

// TODO move route to api
func (s *Handler) Search() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		searchArgs := records.SearchArgs{}
		err := s.formDecoder.Decode(&searchArgs, r.URL.Query())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hits, err := s.service.SearchRecs(searchArgs)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.render.JSON(w, http.StatusOK, hits)
	}

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

func renderSource(j json.RawMessage) template.HTML {
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
					v.Get(`subfields`).ForEach(func(_, subfield gjson.Result) bool {
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
