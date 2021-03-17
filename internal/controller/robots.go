package controller

import (
	"net/http"

	"github.com/ugent-library/momo/engine"
	"github.com/ugent-library/momo/web/render"
)

func Robots(_ engine.Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		render.Text(w, r, "User-agent: *\nDisallow: /\n")
	}
}
