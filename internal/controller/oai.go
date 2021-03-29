package controller

import (
	"net/http"

	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/oaipmh"
)

func OAI(_ engine.Engine) http.Handler {
	return oaipmh.NewProvider()
}
