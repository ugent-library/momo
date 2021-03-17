package controller

import (
	"net/http"

	"github.com/ugent-library/momo/internal/engine"
	"github.com/ugent-library/momo/internal/render"
)

func Robots(_ engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Text(w, r, "User-agent: *\nDisallow: /\n")
	}
}
