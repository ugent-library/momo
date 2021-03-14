package app

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/form/v4"
	"github.com/tidwall/gjson"
	"github.com/ugent-library/momo/engine"
	_ "github.com/ugent-library/momo/themes/orpheus" // load orpheus theme
	_ "github.com/ugent-library/momo/themes/ugent"   // load ugent theme
	"github.com/ugent-library/momo/web/theme"
	"github.com/unrolled/render"
)

var themeCtxKey = &contextKey{"Theme"}
var scopeCtxKey = &contextKey{"Scope"}

type App struct {
	engine.Engine
	Router         chi.Router
	funcMaps       []template.FuncMap
	formDecoder    *form.Decoder
	themeRenderers map[string]*render.Render
}

func New(e engine.Engine) *App {
	a := &App{
		Engine: e,
		Router: chi.NewRouter(),
		funcMaps: []template.FuncMap{{
			"renderTitle":        renderTitle,
			"renderMetadata":     renderMetadata,
			"renderSourceView":   renderSourceView,
			"renderInternalView": renderInternalView,
		}},
		formDecoder:    form.NewDecoder(),
		themeRenderers: make(map[string]*render.Render),
	}

	for _, t := range theme.Themes() {
		r := render.New(render.Options{
			Directory: fmt.Sprintf("themes/%s/templates", t.Name()),
			Layout:    "layout",
			Funcs:     append(a.funcMaps, t.FuncMaps()...),
		})
		a.themeRenderers[t.Name()] = r
	}

	return a
}

func (a *App) DecodeForm(v interface{}, values url.Values) error {
	return a.formDecoder.Decode(v, values)
}

func (a *App) RenderHTML(w http.ResponseWriter, r *http.Request, status int, tmpl string, v interface{}) error {
	t := GetTheme(r)
	renderer, ok := a.themeRenderers[t]
	if !ok {
		panic("momo: theme " + t + " not found")
	}
	return renderer.HTML(w, status, tmpl, v)
}

// taken from chi render
func (a *App) RenderJSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	// if status, ok := r.Context().Value(StatusCtxKey).(int); ok {
	// 	w.WriteHeader(status)
	// }
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (c *contextKey) String() string {
	return c.name
}

// GetTheme gets the theme name from the request context.
func GetTheme(r *http.Request) string {
	return r.Context().Value(themeCtxKey).(string)
}

// SetTheme is a middleware that forces the theme name.
func SetTheme(theme string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), themeCtxKey, theme))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func GetScope(r *http.Request) engine.Scope {
	return r.Context().Value(scopeCtxKey).(engine.Scope)
}

// SetTheme is a middleware that forces the theme name.
func SetScope(scope engine.Scope) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), scopeCtxKey, scope))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// TODO move and cleanup helpers

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
