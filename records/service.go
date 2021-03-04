package records

import "sync"

type Storage interface {
	AddRec(*Rec) error
	GetRec(string) (*Rec, error)
}

type SearchStorage interface {
	CreateIndex() error
	DeleteIndex() error
	AddRecs(<-chan *Rec)
	SearchRecs(SearchArgs) (*Hits, error)
}

type Service interface {
	AddRec(*Rec) error
	AddRecs(<-chan *Rec)
	GetRec(string) (*Rec, error)
	SearchRecs(SearchArgs) (*Hits, error)
}

type service struct {
	store       Storage
	searchStore SearchStorage
	scope       Scope
}

func NewService(store Storage, searchStore SearchStorage, scopes ...Scope) Service {
	var scope Scope
	if len(scopes) > 0 {
		scope = scopes[0]
	}
	return &service{store, searchStore, scope}
}

func (s *service) AddRec(rec *Rec) error {
	if err := s.store.AddRec(rec); err != nil {
		return err
	}
	return nil
}

func (s *service) AddRecs(in <-chan *Rec) {
	var wg sync.WaitGroup

	addRecs := func(in <-chan *Rec, out chan<- *Rec, wg *sync.WaitGroup) {
		defer wg.Done()

		for r := range in {
			s.store.AddRec(r)
			// index after storing
			out <- r
		}
	}

	// indexing channel
	out := make(chan *Rec)

	// start bulk indexer
	go s.searchStore.AddRecs(out)

	// store recs
	for i := 0; i < 4; i++ { // TODO make configurable
		wg.Add(1)
		go addRecs(in, out, &wg)
	}

	// close indexing channel when all recs are stored
	go func() {
		wg.Wait()
		close(out)
	}()
}

func (s *service) GetRec(id string) (*Rec, error) {
	return s.store.GetRec(id)
}

func (s *service) SearchRecs(args SearchArgs) (*Hits, error) {
	return s.searchStore.SearchRecs(args.WithScope(s.scope))
}
