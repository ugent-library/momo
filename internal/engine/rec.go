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

func (e *engine) GetRec(id string) (*Rec, error) {
	return e.store.GetRec(id)
}

func (e *engine) GetAllRecs(c chan<- *Rec) error {
	return e.store.GetAllRecs(c)
}

func (e *engine) SearchRecs(args SearchArgs) (*RecHits, error) {
	return e.searchStore.SearchRecs(args)
}

// TODO don't die
// TODO send errors back over a channel
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
	for i := 0; i < runtime.NumCPU()/2; i++ { // TODO make configurable
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
