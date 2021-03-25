package gettext

import (
	"github.com/leonelquinteros/gotext"
	"github.com/ugent-library/momo/internal/engine"
	"golang.org/x/text/language"
)

type i18n struct {
	langs   []language.Tag
	locales []engine.Locale
	matcher language.Matcher
}

type locale struct {
	*gotext.Locale
	lang language.Tag
}

func New() *i18n {
	langs := []language.Tag{
		language.English, // first is the fallback language
		language.Dutch,
	}
	locales := make([]engine.Locale, len(langs))
	for i, lang := range langs {
		loc := gotext.NewLocale("etc/locales", lang.String())
		loc.AddDomain("default")
		locales[i] = &locale{loc, lang}
	}

	return &i18n{
		langs:   langs,
		locales: locales,
		matcher: language.NewMatcher(langs),
	}

}

func (i *i18n) Locales() []engine.Locale {
	return i.locales
}

func (i *i18n) GetLocale(langs ...string) engine.Locale {
	_, match := language.MatchStrings(i.matcher, langs...)
	return i.locales[match]
}

func (l *locale) Language() language.Tag {
	return l.lang
}
