package records

type Storage interface {
	AddRec(*Rec) error
}

type Service interface {
	AddRec(*Rec) error
}

type service struct {
	store       Storage
	searchStore Storage
}

func NewService(store Storage, searchStore Storage) Service {
	return &service{store: store, searchStore: searchStore}
}

// TODO decouple indexing
func (s *service) AddRec(rec *Rec) error {
	if err := s.store.AddRec(rec); err != nil {
		return err
	}
	if err := s.searchStore.AddRec(rec); err != nil {
		return err
	}
	return nil
}
