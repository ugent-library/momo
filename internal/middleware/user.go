package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/ugent-library/momo/internal/ctx"
)

func SetUser() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("user")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			c := context.WithValue(r.Context(), ctx.UserKey, cookie.Value)
			next.ServeHTTP(w, r.WithContext(c))
		})
	}
}

func RequireUser(authURL string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !ctx.HasUser(r) {
				http.Redirect(w, r, authURL, http.StatusFound)
				return
			}
			log.Printf("%s", ctx.GetUser(r))
			next.ServeHTTP(w, r)
		})
	}
}
