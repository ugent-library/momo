package records

type Storage interface {
	AddRec(*Rec) error
	AddRecs(<-chan *Rec)
}

type Service interface {
	AddRec(*Rec) error
	AddRecs(<-chan *Rec)
}

type service struct {
	store Storage
}

func NewService(store Storage) Service {
	return &service{store: store}
}

func (s *service) AddRec(rec *Rec) error {
	if err := s.store.AddRec(rec); err != nil {
		return err
	}
	return nil
}

func (s *service) AddRecs(c <-chan *Rec) {
	s.store.AddRecs(c)
}
