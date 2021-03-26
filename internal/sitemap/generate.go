package sitemap

import (
	"fmt"
	"time"

	"github.com/ikeikeikeike/go-sitemap-generator/v2/stm"
	"github.com/ugent-library/momo/internal/engine"
)

func Generate(e engine.Engine, dir string) (err error) {
	now := time.Now().Unix()
	sitemaps := make([]string, 0)
	i := 0

	newSm := func() *stm.Sitemap {
		file := fmt.Sprintf("sitemap-%d-%d.xml", now, len(sitemaps))
		sitemaps = append(sitemaps, file)
		sm := stm.NewSitemap(1)
		sm.SetDefaultHost("http://localhost:3000")
		sm.SetSitemapsPath(dir)
		sm.SetFilename(file)
		sm.Create()
		return sm
	}

	sm := newSm()

	c := e.AllRecs()
	defer c.Close()

	for c.Next() {
		if err = c.Error(); err != nil {
			break
		}

		i++

		rec := c.Value()

		sm.Add(stm.URL{
			{"loc", fmt.Sprintf("/collection/%s/%s", rec.Collection, rec.ID)},
		})

		if i == 10000 {
			sm.Finalize()
			break
		}
	}

	return
}

// type sitemap struct {
// 	p string
// 	f *os.File
// }

// func writeSitemap(e engine.Engine, p string, recs []*engine.Rec) (err error) {
// 	f, err := os.Create(p)
// 	if err != nil {
// 		return
// 	}

// 	defer f.Close()

// 	f.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
// 	f.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xhtml="http://www.w3.org/1999/xhtml">`)

// 	f.WriteString(`</urlset>`)

// 	return
// }
