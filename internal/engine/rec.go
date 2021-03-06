package engine

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"
	"time"
)

type RecEngine interface {
	GetRec(string) (*Rec, error)
	EachRec(func(*Rec) bool) error
	SearchRecs(SearchArgs) (*RecHits, error)
	SearchMoreRecs(string) (*RecHits, error)
	SearchEachRec(SearchArgs, func(*Rec) bool) error
	AddRecsBySourceID(<-chan *Rec)
	UpdateRecMetadata(string, map[string]interface{}) (*Rec, error)
	IndexRecs() error
	CreateRecIndex() error
	DeleteRecIndex() error
}

type Rec struct {
	ID             string                 `json:"id"`
	Collection     string                 `json:"collection"`
	Type           string                 `json:"type"`
	Metadata       map[string]interface{} `json:"metadata"`
	SourceID       string                 `json:"sourceID"`
	SourceFormat   string                 `json:"sourceFormat,omitempty"`
	SourceMetadata []byte                 `json:"sourceMetadata,omitempty"`
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

func (e *engine) AddRecsBySourceID(storeC <-chan *Rec) {
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
				err := e.store.AddRecBySourceID(r)
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

func (e *engine) UpdateRecMetadata(id string, m map[string]interface{}) (*Rec, error) {
	rec, err := e.store.UpdateRecMetadata(id, m)
	if err != nil {
		return nil, err
	}
	if err = e.searchStore.AddRec(rec); err != nil {
		return nil, err
	}
	return rec, nil
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
