package metadata

import (
	"github.com/ugent-library/momo/internal/engine"
	"golang.org/x/text/language"
)

type Text struct {
	Lang language.Tag
	Text string
}

type Contributor struct {
	Name string
}

type IIIF struct {
	Manifest      string
	MiradorViewer string
}

type Rec struct {
	*engine.Rec
}

func WrapRec(rec *engine.Rec) Rec {
	return Rec{Rec: rec}
}

func (r Rec) Abstract() []Text {
	return r.textSlice("abstract")
}

func (r Rec) Author() []Contributor {
	return r.contributorSlice("author")
}

func (r Rec) DOI() []string {
	return r.stringSlice("doi")
}

func (r Rec) Edition() string {
	return r.string("edition")
}

func (r Rec) IIIFManifest() string {
	val, ok := r.Metadata["iiif"]
	if !ok {
		return ""
	}
	m, ok := val.(map[string]interface{})
	if !ok {
		return ""
	}
	if v, ok := m["manifest"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (r Rec) IIIFViewer(name string) string {
	val, ok := r.Metadata["iiif"]
	if !ok {
		return ""
	}
	m, ok := val.(map[string]interface{})
	if !ok {
		return ""
	}
	if v, ok := m[name+"Viewer"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (r Rec) ISBN() []string {
	return r.stringSlice("isbn")
}

func (r Rec) Note() []Text {
	return r.textSlice("note")
}

func (r Rec) PlaceOfPublication() string {
	return r.string("placeOfPublication")
}

func (r Rec) PublicationDate() string {
	return r.string("publicationDate")
}

func (r Rec) Publisher() string {
	return r.string("publisher")
}

func (r Rec) Tag() []string {
	return r.stringSlice("tag")
}

func (r Rec) Title() string {
	return r.string("title")
}

func (r Rec) string(field string) string {
	if val, ok := r.Metadata[field]; ok {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

func (r Rec) stringSlice(field string) (ss []string) {
	if vals, ok := r.Metadata[field].([]interface{}); ok {
		for _, val := range vals {
			if s, ok := val.(string); ok {
				ss = append(ss, s)
			}
		}
	}
	return
}

func (r Rec) textSlice(field string) (tt []Text) {
	vals, ok := r.Metadata[field].([]interface{})
	if !ok {
		return
	}
	for _, val := range vals {
		m, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		t := Text{}
		if v, ok := m["lang"].(string); ok {
			if l, err := language.Parse(v); err == nil {
				t.Lang = l
			}
		}
		if v, ok := m["text"].(string); ok {
			t.Text = v
		}
		tt = append(tt, t)
	}
	return
}

func (r Rec) contributorSlice(field string) (cc []Contributor) {
	vals, ok := r.Metadata[field].([]interface{})
	if !ok {
		return
	}
	for _, val := range vals {
		m, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		c := Contributor{}
		if v, ok := m["name"].(string); ok {
			c.Name = v
		}
		cc = append(cc, c)
	}
	return
}
