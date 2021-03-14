package engine

type Engine interface {
	RecEngine
	LensEngine
}

type engine struct {
	store       Storage
	searchStore SearchStorage
	lenses      map[string]*Lens
}

type option func(*engine)

func New(opts ...option) Engine {
	e := &engine{}
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
