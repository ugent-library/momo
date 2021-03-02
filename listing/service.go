package listing

type Storage interface {
	SearchRecs(SearchArgs) (*RecHits, error)
}

type Service interface {
	SearchRecs(SearchArgs) (*RecHits, error)
}

type service struct {
	store Storage
	scope Scope
}

func NewService(store Storage, scope Scope) Service {
	return &service{store: store, scope: scope}
}

func (s *service) SearchRecs(args SearchArgs) (*RecHits, error) {
	return s.store.SearchRecs(args.WithScope(s.scope))
}
