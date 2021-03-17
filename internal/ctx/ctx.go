package ctx

import (
	"context"
	"net/http"

	"github.com/ugent-library/momo/engine"
)

var (
	StatusKey = &key{"Status"}
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

func SetStatus(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), StatusKey, status))
}

func GetScope(r *http.Request) engine.Scope {
	return r.Context().Value(ScopeKey).(engine.Scope)
}

// GetTheme gets the theme name from the request context.
func GetTheme(r *http.Request) string {
	return r.Context().Value(ThemeKey).(string)
}
