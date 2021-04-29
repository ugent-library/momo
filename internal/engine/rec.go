package engine

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
)

type RecEngine interface {
	GetRec(string) (*Rec, error)
	EachRec(func(*Rec) bool) error
	SearchRecs(SearchArgs) (*RecHits, error)
	SearchMoreRecs(string) (*RecHits, error)
	SearchEachRec(SearchArgs, func(*Rec) bool) error
	AddRecs(<-chan *Rec)
	IndexRecs() error
	CreateRecIndex() error
	DeleteRecIndex() error
}

type Rec struct {
	ID             string                 `json:"id"`
	Collection     string                 `json:"collection"`
	Type           string                 `json:"type"`
	Metadata       map[string]interface{} `json:"metadata"`
	Source         string                 `json:"source,omitempty"`
	SourceID       string                 `json:"source_id,omitempty"`
	SourceFormat   string                 `json:"source_format,omitempty"`
	SourceMetadata []byte                 `json:"source_metadata,omitempty"`
	CreatedAt      time.Time              `json:"createdAt"`
	UpdatedAt      time.Time              `json:"updatedAt"`
}

type RecHits struct {
	CursorID       string          `json:"-"`
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

func (e *engine) EachRec(fn func(*Rec) bool) error {
	return e.store.EachRec(fn)
}

func (e *engine) SearchRecs(args SearchArgs) (*RecHits, error) {
	return e.searchStore.SearchRecs(args)
}

func (e *engine) SearchMoreRecs(cursorID string) (*RecHits, error) {
	return e.searchStore.SearchMoreRecs(cursorID)
}

func (e *engine) SearchEachRec(args SearchArgs, fn func(*Rec) bool) error {
	if args.Size == 0 {
		args.Size = 200
	}
	args.Cursor = true
	hits, err := e.searchStore.SearchRecs(args)
	if err != nil {
		return err
	}
	for len(hits.Hits) > 0 {
		for _, hit := range hits.Hits {
			if ok := fn(&hit.Rec); !ok {
				return nil
			}
		}
		if hits, err = e.searchStore.SearchMoreRecs(hits.CursorID); err != nil {
			return err
		}
	}
	return nil
}

func (e *engine) AddRecs(storeC <-chan *Rec) {
	var storeWG sync.WaitGroup
	var indexWG sync.WaitGroup

	// indexing channel
	indexC := make(chan *Rec)

	// start bulk indexer
	go func() {
		indexWG.Add(1)
		defer indexWG.Done()
		e.searchStore.AddRecs(indexC)
	}()

	// store recs
	for i := 0; i < runtime.NumCPU()/2; i++ {
		storeWG.Add(1)
		go func() {
			defer storeWG.Done()
			for r := range storeC {
				if r.ID == "" {
					r.ID = uuid.NewString()
				}
				err := e.store.AddRec(r)
				if err != nil {
					log.Fatal(err)
				}
				// index after storing
				indexC <- r
			}
		}()
	}

	// close indexing channel when all recs are stored
	storeWG.Wait()
	close(indexC)
	// wait for indexing to finish
	indexWG.Wait()
}

func (e *engine) IndexRecs() (err error) {
	var indexWG sync.WaitGroup

	// indexing channel
	indexC := make(chan *Rec)

	go func() {
		indexWG.Add(1)
		defer indexWG.Done()
		e.searchStore.AddRecs(indexC)
	}()

	// send recs to indexer
	e.EachRec(func(rec *Rec) bool {
		indexC <- rec
		return true
	})

	close(indexC)

	// wait for indexing to finish
	indexWG.Wait()

	return
}

func (e *engine) CreateRecIndex() error {
	return e.searchStore.CreateRecIndex()
}

func (e *engine) DeleteRecIndex() error {
	return e.searchStore.DeleteRecIndex()
}
