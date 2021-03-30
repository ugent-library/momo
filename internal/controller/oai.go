package controller

import (
	"net/http"

	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/oaipmh"
)

func OAI(_ engine.Engine) http.Handler {
	return oaipmh.NewProvider(oaipmh.ProviderOptions{
		RepositoryName:    "momo",
		BaseURL:           "http://myserver.org",
		AdminEmail:        []string{"momo@myserver.org"},
		EarliestDatestamp: "1970-01-01T00:00:01Z",
		MetadataFormats: []oaipmh.MetadataFormat{
			oaipmh.OAIDC,
		},
		Sets: []oaipmh.Set{},
	})
}
