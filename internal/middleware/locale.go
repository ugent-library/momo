package middleware

import (
	"context"
	"net/http"

	"github.com/ugent-library/momo/internal/ctx"
	"github.com/ugent-library/momo/internal/engine"
)

// SetLocale is a middleware that sets the locale based on the Accept-Language
// header if not already set.
func SetLocale(e engine.Engine) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			// skip if already set
			if r.Context().Value(ctx.LocaleKey) == nil {
				loc := e.GetLocale(r.Header.Get("Accept-Language"))
				c := context.WithValue(r.Context(), ctx.LocaleKey, loc)
				r = r.WithContext(c)
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
