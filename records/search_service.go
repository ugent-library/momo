package records

type SearchStorage interface {
	SearchRecs(SearchArgs) (*Hits, error)
}

type SearchService interface {
	Search(SearchArgs) (*Hits, error)
}

type searchService struct {
	store SearchStorage
	scope Scope
}

func NewSearchService(store SearchStorage, scope Scope) SearchService {
	return &searchService{store: store, scope: scope}
}

func (s *searchService) Search(args SearchArgs) (*Hits, error) {
	return s.store.SearchRecs(args.WithScope(s.scope))
}
