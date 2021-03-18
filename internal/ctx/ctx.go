package ctx

import (
	"context"
	"net/http"

	"github.com/ugent-library/momo/internal/engine"
)

var (
	StatusKey = &key{"Status"}
	LocaleKey = &key{"Locale"}
	ThemeKey  = &key{"Theme"}
	ScopeKey  = &key{"Scope"}
)

type key struct {
	name string
}

func (c *key) String() string {
	return c.name
}

func GetStatus(r *http.Request) int {
	if status, ok := r.Context().Value(StatusKey).(int); ok {
		return status
	}
	return http.StatusOK
}

func SetStatus(r *http.Request, s int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusKey, s))
}

func GetLocale(r *http.Request) engine.Locale {
	return r.Context().Value(LocaleKey).(engine.Locale)
}

func GetScope(r *http.Request) engine.Scope {
	return r.Context().Value(ScopeKey).(engine.Scope)
}

// GetTheme gets the theme name from the request context.
func GetTheme(r *http.Request) string {
	return r.Context().Value(ThemeKey).(string)
}
