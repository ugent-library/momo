package engine

import "sync"

type Engine interface {
	RecEngine
	RecEncoderEngine
	LensEngine
	I18nEngine
	Reset() error
}

type engine struct {
	store         Storage
	searchStore   SearchStorage
	recEncoders   map[string]RecEncoderFactory
	recEncodersMu sync.Mutex
	lenses        []Lens
	I18nEngine
}

type option func(*engine)

func New(opts ...option) Engine {
	e := &engine{
		recEncoders: make(map[string]RecEncoderFactory),
	}

	for _, opt := range opts {
		opt(e)
	}

	e.initLens()

	return e
}

func WithStore(s Storage) option {
	return func(e *engine) {
		e.store = s
	}
}

func WithSearchStore(s SearchStorage) option {
	return func(e *engine) {
		e.searchStore = s
	}
}

func WithRecEncoder(format string, factory RecEncoderFactory) option {
	return func(e *engine) {
		e.recEncoders[format] = factory
	}
}

func WithI18n(i18n I18nEngine) option {
	return func(e *engine) {
		e.I18nEngine = i18n
	}
}

func (e *engine) Reset() error {
	if err := e.searchStore.Reset(); err != nil {
		return err
	}
	return e.store.Reset()
}
