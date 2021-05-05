package ctx

import (
	"context"
	"net/http"

	"github.com/ugent-library/momo/internal/engine"
)

var (
	// CollectionKey is used to get or set the collection in the request context
	CollectionKey = &key{"Collection"}
	// LocaleKey is used to get or set the locale in the request context
	LocaleKey = &key{"Locale"}
	// StatusKey is used to get or set the http status in the request context
	StatusKey = &key{"Status"}
	// ThemeKey is used to get or set the ui theme in the request context
	ThemeKey = &key{"Theme"}
	// UserKey is used to get or set the user in the request context
	UserKey = &key{"User"}
)

type key struct {
	name string
}

func (c *key) String() string {
	return c.name
}

// GetCollection gets the current scope from the request context.
func GetCollection(r *http.Request) string {
	return r.Context().Value(CollectionKey).(string)
}

// GetLocale gets the current locale from the request context.
func GetLocale(r *http.Request) engine.Locale {
	if v, ok := r.Context().Value(LocaleKey).(engine.Locale); ok {
		return v
	}
	return nil
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

// GetTheme gets the current theme name from the request context.
func GetTheme(r *http.Request) string {
	return r.Context().Value(ThemeKey).(string)
}

func HasUser(r *http.Request) bool {
	_, ok := r.Context().Value(UserKey).(string)
	return ok
}

func GetUser(r *http.Request) string {
	return r.Context().Value(UserKey).(string)
}
