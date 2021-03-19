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
	if status, ok := r.Context().Value(StatusKey).(int); ok {
		return status
	}
	return http.StatusOK
}

func SetStatus(r *http.Request, s int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusKey, s))
}

// TODO return default
func GetLocale(r *http.Request) engine.Locale {
	return r.Context().Value(LocaleKey).(engine.Locale)
}

// TODO return default
func GetScope(r *http.Request) engine.Scope {
	return r.Context().Value(ScopeKey).(engine.Scope)
}

// GetTheme gets the theme name from the request context.
// TODO return default
func GetTheme(r *http.Request) string {
	return r.Context().Value(ThemeKey).(string)
}
