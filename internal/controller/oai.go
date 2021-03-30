package controller

import (
	"bytes"
	"net/http"
	"strings"
	"time"

	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/oaipmh"
)

func OAI(e engine.Engine) http.Handler {
	return oaipmh.NewProvider(oaipmh.ProviderOptions{
		RepositoryName:    "momo",
		BaseURL:           "http://myserver.org",
		AdminEmail:        []string{"momo@myserver.org"},
		EarliestDatestamp: "1970-01-01T00:00:01Z",
		MetadataFormats: []oaipmh.MetadataFormat{
			oaipmh.OAIDC,
		},
		Sets: []oaipmh.Set{},

		GetRecord: func(id string, fmt string) *oaipmh.Record {
			parts := strings.Split(id, ":")
			rec, _ := e.GetRec(parts[0], strings.Join(parts[1:], ":"))
			if rec == nil {
				return nil
			}
			var b bytes.Buffer
			enc := e.NewRecEncoder(&b, "oai_dc")
			enc.Encode(rec)
			r := oaipmh.Record{
				// TODO Sets
				Header: oaipmh.Header{
					Datestamp:  rec.UpdatedAt.UTC().Format(time.RFC3339),
					Identifier: id,
				},
				Metadata: oaipmh.Metadata{XML: b.Bytes()},
			}
			return &r
		},

		ListRecords: func() []*oaipmh.Record {
			hits, _ := e.SearchRecs(engine.SearchArgs{Size: 100})
			var records []*oaipmh.Record
			for _, hit := range hits.Hits {
				var b bytes.Buffer
				enc := e.NewRecEncoder(&b, "oai_dc")
				enc.Encode(&hit.Rec)
				r := oaipmh.Record{
					// TODO Sets
					Header: oaipmh.Header{
						Datestamp:  hit.UpdatedAt.UTC().Format(time.RFC3339),
						Identifier: strings.Join([]string{hit.Collection, hit.ID}, ":"),
					},
					Metadata: oaipmh.Metadata{XML: b.Bytes()},
				}
				records = append(records, &r)
			}
			return records
		},
	})
}
