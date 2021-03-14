package engine

import (
	"encoding/json"
	"sync"
	"time"
)

type RecEngine interface {
	GetRec(string) (*Rec, error)
	GetAllRecs(chan<- *Rec) error
	SearchRecs(SearchArgs) (*RecHits, error)
	AddRecs(<-chan *Rec)
	IndexRecs() error
	CreateRecIndex() error
	DeleteRecIndex() error
}

type Rec struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Collection  []string        `json:"collection"`
	RawMetadata json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	RawSource   json.RawMessage `json:"source"`
}

type RecHits struct {
	Total          int             `json:"total"`
	Hits           []*RecHit       `json:"hits"`
	RawAggregation json.RawMessage `json:"aggregation"`
}

type RecHit struct {
	Rec
	RawHighlight json.RawMessage `json:"highlight"`
}

func (e *engine) GetRec(id string) (*Rec, error) {
	return e.store.GetRec(id)
}

func (e *engine) GetAllRecs(c chan<- *Rec) error {
	return e.store.GetAllRecs(c)
}

func (e *engine) SearchRecs(args SearchArgs) (*RecHits, error) {
	return e.searchStore.SearchRecs(args)
}

func (e *engine) AddRecs(in <-chan *Rec) {
	var wg sync.WaitGroup

	addRecs := func(in <-chan *Rec, out chan<- *Rec, wg *sync.WaitGroup) {
		defer wg.Done()

		for r := range in {
			e.store.AddRec(r)
			// index after storing
			out <- r
		}
	}

	// indexing channel
	out := make(chan *Rec)

	// start bulk indexer
	go e.searchStore.AddRecs(out)

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

func (e *engine) IndexRecs() error {
	c := make(chan *Rec)
	defer close(c)
	go e.searchStore.AddRecs(c)
	err := e.store.GetAllRecs(c)
	return err
}

func (e *engine) CreateRecIndex() error {
	return e.searchStore.CreateRecIndex()
}

func (e *engine) DeleteRecIndex() error {
	return e.searchStore.DeleteRecIndex()
}
