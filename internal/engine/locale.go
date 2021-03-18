package engine

import (
	"github.com/leonelquinteros/gotext"
	"golang.org/x/text/language"
)

var (
	languages = []language.Tag{
		language.English, // first is the fallback language
		language.Dutch,
	}

	languageMatcher = language.NewMatcher(languages)

	locales []Locale
)

type LocaleEngine interface {
	Languages() []language.Tag
	GetLocale(...string) Locale
}

type Locale interface {
	Get(string, ...interface{}) string
}

func init() {
	for _, tag := range languages {
		lang := tag.String()
		locale := gotext.NewLocale("etc/locales", lang)
		locale.AddDomain("default")
		locales = append(locales, locale)
	}
}

func (e *engine) Languages() []language.Tag {
	return languages
}

func (e *engine) GetLocale(langs ...string) Locale {
	_, i := language.MatchStrings(languageMatcher, langs...)
	return locales[i]
}
