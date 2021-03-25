package engine

import "golang.org/x/text/language"

type I18nEngine interface {
	Locales() []Locale
	DefaultLocale() Locale
	GetLocale(...string) Locale
}

type Locale interface {
	Language() language.Tag
	Get(string, ...interface{}) string
}
