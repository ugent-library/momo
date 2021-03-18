package middleware

import (
	"context"
	"net/http"

	"github.com/ugent-library/momo/internal/ctx"
	"github.com/ugent-library/momo/internal/engine"
)

func SetLocale(e engine.Engine) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			loc := e.GetLocale(
				r.URL.Query().Get("lang"),
				r.Header.Get("Accept-Language"),
			)
			r = r.WithContext(context.WithValue(r.Context(), ctx.LocaleKey, loc))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// SetTheme is a middleware that forces the theme name.
func SetTheme(theme string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ctx.ThemeKey, theme))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// SetScope is a middleware that forces the scope.
func SetScope(scope engine.Scope) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ctx.ScopeKey, scope))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
