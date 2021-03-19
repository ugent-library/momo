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
	Locales() []Locale
	GetLocale(...string) Locale
}

type Locale interface {
	Language() language.Tag
	Get(string, ...interface{}) string
}

type locale struct {
	*gotext.Locale
	lang language.Tag
}

func init() {
	for _, lang := range languages {
		loc := gotext.NewLocale("etc/locales", lang.String())
		loc.AddDomain("default")
		locales = append(locales, &locale{loc, lang})
	}
}

func (e *engine) Locales() []Locale {
	return locales
}

func (e *engine) GetLocale(langs ...string) Locale {
	_, match := language.MatchStrings(languageMatcher, langs...)
	return locales[match]
}

func (l *locale) Language() language.Tag {
	return l.lang
}
