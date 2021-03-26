package sitemap

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path"

	"github.com/ugent-library/momo/internal/engine"
)

func Generate(e engine.Engine, baseURL string) (err error) {
	var langs []string
	for _, loc := range e.Locales() {
		langs = append(langs, loc.Language().String())
	}

	sm := &sitemap{
		path: "static/sitemaps",
	}

	c := e.AllRecs()
	defer c.Close()

	for c.Next() {
		if err = c.Error(); err != nil {
			break
		}

		rec := c.Value()

		for _, lang := range langs {
			alts := make([]alternate, len(langs))
			for i, l := range langs {
				alts[i].lang = l
				alts[i].url = fmt.Sprintf("%s/%s/collection/%s/%s", baseURL, l, rec.Collection, rec.ID)
			}
			alts = append(alts, alternate{
				lang: "x-default",
				url:  fmt.Sprintf("%s/collection/%s/%s", baseURL, rec.Collection, rec.ID),
			})
			sm.add(sitemapURL{
				loc:        fmt.Sprintf("%s/%s/collection/%s/%s", baseURL, lang, rec.Collection, rec.ID),
				lastmod:    rec.UpdatedAt.Format("2006-01-02"),
				priority:   "0.9",
				alternates: alts,
			})
		}
	}

	sm.finish()

	return
}

type sitemap struct {
	path     string
	numFiles int
	numURLs  int
	f        *os.File
	w        *bufio.Writer
}

type sitemapURL struct {
	loc        string
	lastmod    string
	priority   string
	alternates []alternate
}

type alternate struct {
	lang string
	url  string
}

func (sm *sitemap) switchFile() error {
	sm.finish()

	if err := os.MkdirAll(sm.path, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path.Join(sm.path, fmt.Sprintf("sitemap-%d.xml", sm.numFiles)))
	if err != nil {
		return err
	}

	sm.f = f
	sm.w = bufio.NewWriter(f)

	sm.w.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sm.w.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">`)

	sm.numFiles++

	return nil
}

func (sm *sitemap) add(u sitemapURL) error {
	if sm.numURLs == 0 || sm.numURLs == 50000 {
		if err := sm.switchFile(); err != nil {
			return err
		}
		sm.numURLs = 0
	}

	w := sm.w
	w.WriteString(`<url>`)
	w.WriteString(`<loc>`)
	xml.EscapeText(w, []byte(u.loc))
	w.WriteString(`</loc>`)
	for _, alt := range u.alternates {
		w.WriteString(`<xhtml:link rel="alternate" hreflang="`)
		xml.EscapeText(w, []byte(alt.lang))
		w.WriteString(`" href="`)
		xml.EscapeText(w, []byte(alt.url))
		w.WriteString(`"/>`)
	}
	w.WriteString(`</url>`)

	sm.numURLs++

	return nil
}

func (sm *sitemap) finish() {
	if sm.f != nil {
		sm.w.WriteString(`</urlset>`)

		sm.w.Flush()
		sm.f.Close()
	}
}
