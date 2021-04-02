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
		Sets: []oaipmh.Set{
			{SetName: "all", SetSpec: "All records"},
		},

		GetRecord: func(r *oaipmh.Request) *oaipmh.Record {
			parts := strings.Split(r.Identifier, ":")
			rec, _ := e.GetRec(parts[0], strings.Join(parts[1:], ":"))
			if rec == nil {
				return nil
			}
			var b bytes.Buffer
			enc := e.NewRecEncoder(&b, "oai_dc")
			enc.Encode(rec)
			oaiRec := oaipmh.Record{
				// TODO Sets
				Header: &oaipmh.Header{
					Datestamp:  rec.UpdatedAt.UTC().Format(time.RFC3339),
					Identifier: r.Identifier,
				},
				Metadata: &oaipmh.Payload{XML: b.Bytes()},
			}
			return &oaiRec
		},

		// TODO from, until, set, metadataPrefix
		ListIdentifiers: func(r *oaipmh.Request) ([]*oaipmh.Header, *oaipmh.ResumptionToken) {
			var hits *engine.RecHits
			var err error

			if r.ResumptionToken == "" {
				hits, err = e.SearchRecs(engine.SearchArgs{Size: 200, Cursor: true})
			} else {
				hits, err = e.SearchMoreRecs(r.ResumptionToken)
			}

			if err != nil {
				panic(err) // TODO don't die
			}

			var oaiHeaders []*oaipmh.Header
			for _, hit := range hits.Hits {
				var b bytes.Buffer
				enc := e.NewRecEncoder(&b, "oai_dc")
				enc.Encode(&hit.Rec)
				oaiHeader := oaipmh.Header{
					Datestamp:  hit.UpdatedAt.UTC().Format(time.RFC3339),
					Identifier: strings.Join([]string{hit.Collection, hit.ID}, ":"),
				}
				oaiHeaders = append(oaiHeaders, &oaiHeader)
			}

			token := &oaipmh.ResumptionToken{
				ExpirationDate:   time.Now().Add(time.Minute).UTC().Format(time.RFC3339),
				CompleteListSize: hits.Total,
				Value:            hits.CursorID,
			}

			return oaiHeaders, token
		},

		// TODO from, until, set, metadataPrefix
		ListRecords: func(r *oaipmh.Request) ([]*oaipmh.Record, *oaipmh.ResumptionToken) {
			var hits *engine.RecHits
			var err error

			if r.ResumptionToken == "" {
				hits, err = e.SearchRecs(engine.SearchArgs{Size: 200, Cursor: true})
			} else {
				hits, err = e.SearchMoreRecs(r.ResumptionToken)
			}

			if err != nil {
				panic(err) // TODO don't die
			}

			var records []*oaipmh.Record
			for _, hit := range hits.Hits {
				var b bytes.Buffer
				enc := e.NewRecEncoder(&b, "oai_dc")
				enc.Encode(&hit.Rec)
				r := oaipmh.Record{
					// TODO Sets
					Header: &oaipmh.Header{
						Datestamp:  hit.UpdatedAt.UTC().Format(time.RFC3339),
						Identifier: strings.Join([]string{hit.Collection, hit.ID}, ":"),
					},
					Metadata: &oaipmh.Payload{XML: b.Bytes()},
				}
				records = append(records, &r)
			}

			token := &oaipmh.ResumptionToken{
				ExpirationDate:   time.Now().Add(time.Minute).UTC().Format(time.RFC3339),
				CompleteListSize: hits.Total,
				Value:            hits.CursorID,
			}

			return records, token
		},
	})
}
