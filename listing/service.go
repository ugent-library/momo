package listing

type Storage interface {
	SearchRecs(string, ...map[string]interface{}) (*RecHits, error)
}

type Service interface {
	SearchRecs(string, ...map[string]interface{}) (*RecHits, error)
}

type service struct {
	store Storage
}

func NewService(s Storage) Service {
	return &service{s}
}

func (s *service) SearchRecs(q string, _ ...map[string]interface{}) (*RecHits, error) {
	return s.store.SearchRecs(q)
}
