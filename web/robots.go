package web

import (
	"net/http"

	"github.com/ugent-library/momo/web/app"
)

func Robots(a *app.App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.RenderText(w, r, "User-agent: *\nDisallow: /\n")
	}
}
