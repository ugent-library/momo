package listing

type Storage interface {
	SearchRecs(string) (*RecHits, error)
}

type Service interface {
	SearchRecs(string) (*RecHits, error)
}

type service struct {
	store Storage
}

func NewService(s Storage) Service {
	return &service{s}
}

func (s *service) SearchRecs(q string) (*RecHits, error) {
	return s.store.SearchRecs(q)
}
