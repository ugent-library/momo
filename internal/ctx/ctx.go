package ctx

import (
	"context"
	"net/http"

	"github.com/ugent-library/momo/internal/engine"
)

var (
	// StatusKey is used to get or set the http status in the request context
	StatusKey = &key{"Status"}
	// LocaleKey is used to get or set the locale in the request context
	LocaleKey = &key{"Locale"}
	// ThemeKey is used to get or set the ui theme in the request context
	ThemeKey = &key{"Theme"}
	// ScopeKey is used to get or set the search scope in the request context
	ScopeKey = &key{"Scope"}
)

type key struct {
	name string
}

func (c *key) String() string {
	return c.name
}

// GetStatus gets the http status from the request context if set or 200 OK if not.
func GetStatus(r *http.Request) int {
	if v, ok := r.Context().Value(StatusKey).(int); ok {
		return v
	}
	return http.StatusOK
}

// SetStatus sets the http status in the request context.
func SetStatus(r *http.Request, s int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusKey, s))
}

// GetLocale gets the current locale from the request context.
func GetLocale(r *http.Request) engine.Locale {
	if v, ok := r.Context().Value(LocaleKey).(engine.Locale); ok {
		return v
	}
	return nil
}

// GetScope gets the current scope from the request context.
func GetScope(r *http.Request) engine.Scope {
	if v, ok := r.Context().Value(ScopeKey).(engine.Scope); ok {
		return v
	}
	return nil
}

// GetTheme gets the current theme name from the request context.
func GetTheme(r *http.Request) string {
	return r.Context().Value(ThemeKey).(string)
}
