package sitemap

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/ugent-library/momo/internal/engine"
)

func Generate(e engine.Engine, baseURL string) (err error) {
	var langs []string
	for _, loc := range e.Locales() {
		langs = append(langs, loc.Language().String())
	}

	sm := &sitemap{
		sitemapURL: baseURL + "/s/sitemaps",
		path:       "static/sitemaps",
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

			err = sm.add(sitemapURL{
				loc:        fmt.Sprintf("%s/%s/collection/%s/%s", baseURL, lang, rec.Collection, rec.ID),
				lastmod:    rec.UpdatedAt.Format("2006-01-02"),
				priority:   "0.9",
				alternates: alts,
			})
			if err != nil {
				return
			}
		}
	}

	err = sm.finish()

	return
}

type sitemap struct {
	sitemapURL string
	path       string
	numFiles   int
	numURLs    int
	f          *os.File
	w          *bufio.Writer
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

func (sm *sitemap) switchFile(name string) error {
	sm.closeFile()

	if err := os.MkdirAll(sm.path, os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path.Join(sm.path, name))
	if err != nil {
		return err
	}

	sm.f = f
	sm.w = bufio.NewWriter(f)

	return nil
}

func (sm *sitemap) closeFile() {
	if sm.f != nil {
		sm.w.Flush()
		sm.f.Close()
	}
}

func (sm *sitemap) switchSitemap() error {
	if sm.numFiles > 0 {
		sm.w.WriteString(`</urlset>`)
	}

	if err := sm.switchFile(fmt.Sprintf("sitemap-%d.xml", sm.numFiles)); err != nil {
		return err
	}

	sm.w.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	sm.w.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">`)

	sm.numFiles++
	sm.numURLs = 0

	return nil
}

func (sm *sitemap) add(u sitemapURL) error {
	if sm.numURLs == 0 || sm.numURLs == 50000 {
		if err := sm.switchSitemap(); err != nil {
			return err
		}
	}

	w := sm.w
	w.WriteString(`<url><loc>`)
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

func (sm *sitemap) finish() error {
	if err := sm.switchFile("sitemap.xml"); err != nil {
		return err
	}

	w := sm.w

	w.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	w.WriteString(`<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)

	lastmod := time.Now().Format("2006-01-02")

	for i := 0; i < sm.numFiles; i++ {
		url := path.Join(sm.sitemapURL, fmt.Sprintf("sitemap-%d.xml", i))
		w.WriteString(`<sitemap><loc>`)
		xml.EscapeText(w, []byte(url))
		w.WriteString(`</loc><lastmod>`)
		xml.EscapeText(w, []byte(lastmod))
		w.WriteString(`</lastmod></sitemap>`)
	}

	w.WriteString(`</sitemapindex>`)

	sm.closeFile()

	return nil
}
