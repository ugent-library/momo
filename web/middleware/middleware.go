package middleware

import (
	"context"
	"net/http"

	"github.com/ugent-library/momo/engine"
	"github.com/ugent-library/momo/web/ctx"
)

// ThemeSetter is a middleware that forces the theme name.
func ThemeSetter(theme string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ctx.ThemeKey, theme))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// ScopeSetter is a middleware that forces the scope.
func ScopeSetter(scope engine.Scope) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(context.WithValue(r.Context(), ctx.ScopeKey, scope))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
