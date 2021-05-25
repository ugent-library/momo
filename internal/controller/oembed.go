package controller

import (
	"bytes"
	"encoding/xml"
	"html"
	"log"
	"net/http"

	"github.com/ugent-library/momo/internal/form"
	"github.com/ugent-library/momo/internal/render"
)

type oEmbedRequest struct {
	URL       string `form:"url"`
	Format    string `form:"format"`
	MaxWidth  string `form:"maxwidth"`
	MaxHeight string `form:"maxheight"`
}

type oEmbedResponse struct {
	XMLName         xml.Name `json:"-" xml:"oembed"`
	Type            string   `json:"type,omitempty" xml:"type,omitempty"`
	Version         string   `json:"version,omitempty" xml:"version,omitempty"`
	Title           string   `json:"title,omitempty" xml:"title,omitempty"`
	AuthorName      string   `json:"author_name,omitempty" xml:"author_name,omitempty"`
	AuthorUrl       string   `json:"author_url,omitempty" xml:"author_url,omitempty"`
	ProviderName    string   `json:"provider_name,omitempty" xml:"provider_name,omitempty"`
	ProviderUrl     string   `json:"provider_url,omitempty" xml:"provider_url,omitempty"`
	CacheAge        uint64   `json:"cache_age,omitempty" xml:"cache_age,omitempty"`
	ThumbnailUrl    string   `json:"thumbnail_url,omitempty" xml:"thumbnail_url,omitempty"`
	ThumbnailWidth  int      `json:"thumbnail_width,omitempty" xml:"thumbnail_width,omitempty"`
	ThumbnailHeight int      `json:"thumbnail_height,omitempty" xml:"thumbnail_height,omitempty"`
	URL             string   `json:"url,omitempty" xml:"url,omitempty"`
	HTML            string   `json:"html,omitempty" xml:"html,omitempty"`
	Width           int      `json:"width,omitempty" xml:"width,omitempty"`
	Height          int      `json:"height,omitempty" xml:"height,omitempty"`
}

func OEmbed() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := oEmbedRequest{}
		if err := form.Decode(&req, r.URL.Query()); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		iframe := `<iframe src="` + html.EscapeString(req.URL+`/viewer`) + `" width="480" height="320"></iframe>`

		if req.Format == "xml" {
			var b bytes.Buffer
			xml.EscapeText(&b, []byte(iframe))
			iframe = b.String()
		}

		res := oEmbedResponse{
			Type:    "rich",
			Version: "1.0",
			HTML:    iframe,
		}

		if req.Format == "xml" {
			render.XML(w, r, res)
		} else {
			render.JSON(w, r, res)
		}
	}
}
