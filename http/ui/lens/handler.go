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
	renderDtDd(&b, &m, "author.#.name", "Authors")
	renderDtDd(&b, &m, "abstract.#.text", "Abstract")
	renderDtDd(&b, &m, "edition", "Edition")
	renderDtDd(&b, &m, "publisher", "Publisher")
	renderDtDd(&b, &m, "placeOfPublication", "Place of publication")
	renderDtDd(&b, &m, "publicationDate", "Date published")
	renderDtDd(&b, &m, "doi", "DOI")
	renderDtDd(&b, &m, "isbn", "ISBN")
	renderDtDd(&b, &m, "note.#.text", "Note")
	b.WriteString("</dl>")
	return template.HTML(b.String())
}

func renderDtDd(b *strings.Builder, m *gjson.Result, path string, dt string) {
	res := m.Get(path)
	if res.Exists() {
		b.WriteString("<dt>" + html.EscapeString(dt) + "</dt>")
		res.ForEach(func(_, v gjson.Result) bool {
			b.WriteString("<dd>" + html.EscapeString(v.String()) + "</dd>")
			return true
		})
	}
}
