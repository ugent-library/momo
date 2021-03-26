package engine

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"
	"time"
)

type RecEngine interface {
	GetRec(string, string) (*Rec, error)
	AllRecs() RecCursor
	SearchRecs(SearchArgs) (*RecHits, error)
	AddRecs(<-chan *Rec)
	IndexRecs() error
	CreateRecIndex() error
	DeleteRecIndex() error
}

type RecCursor interface {
	Next() bool
	Value() *Rec
	Error() error
	Close()
}

type Rec struct {
	ID          string          `json:"id"`
	Collection  string          `json:"collection"`
	Type        string          `json:"type"`
	RawMetadata json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	RawSource   json.RawMessage `json:"source"`
	metadata    *RecMetadata
}

type RecMetadata struct {
	Title              string
	Author             []Contributor
	Abstract           []Text
	Edition            string
	Publisher          string
	PlaceOfPublication string
	PublicationDate    string
	DOI                []string
	ISBN               []string
	Note               []Text
}

type Contributor struct {
	Name string
}

type Text struct {
	Lang string
	Text string
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

func (rec *Rec) Metadata() *RecMetadata {
	if rec.metadata == nil {
		rec.metadata = &RecMetadata{}
		if err := json.Unmarshal(rec.RawMetadata, rec.metadata); err != nil {
			panic("momo: invalid metadata in rec " + rec.ID + ": " + err.Error())
		}
	}
	return rec.metadata
}

func (e *engine) GetRec(coll string, id string) (*Rec, error) {
	return e.store.GetRec(coll, id)
}

func (e *engine) AllRecs() RecCursor {
	return e.store.AllRecs()
}

func (e *engine) SearchRecs(args SearchArgs) (*RecHits, error) {
	return e.searchStore.SearchRecs(args)
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
	c := e.store.AllRecs()
	defer c.Close()
	for c.Next() {
		if err = c.Error(); err != nil {
			break
		}
		indexC <- c.Value()
	}

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
