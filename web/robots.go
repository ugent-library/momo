package web

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ugent-library/momo/web/app"
)

func GetRobots(_ *app.App, r chi.Router) {
	r.Get("/robots.txt", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "User-agent: *\nDisallow: /\n")
	})
}
