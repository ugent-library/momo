package engine

import (
	"encoding/json"
	"log"
	"runtime"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"
	"github.com/jmespath/go-jmespath"
)

var (
	getterCache *lru.Cache
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

	AddRepresentation(string, string, []byte) error
}

type Rec struct {
	ID          string          `json:"id"`
	Collection  string          `json:"collection"`
	Type        string          `json:"type"`
	RawMetadata json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"createdAt"`
	UpdatedAt   time.Time       `json:"updatedAt"`
	RawSource   json.RawMessage `json:"source"`
	metadata    map[string]interface{}
}

// type Contributor struct {
// 	Name string
// }

// type Text struct {
// 	Lang string
// 	Text string
// }

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

func init() {
	var err error
	getterCache, err = lru.New(512)
	if err != nil {
		log.Panic(err)
	}
}

func (r *Rec) parseMetadata() {
	if r.metadata != nil {
		return
	}
	r.metadata = make(map[string]interface{})
	if err := json.Unmarshal(r.RawMetadata, &r.metadata); err != nil {
		panic(err)
	}
}

func (r *Rec) Get(path string) interface{} {
	r.parseMetadata()
	var jp *jmespath.JMESPath
	if c, ok := getterCache.Get(path); ok {
		jp = c.(*jmespath.JMESPath)
	} else {
		jp = jmespath.MustCompile(path)
		getterCache.Add(path, jp)
	}
	v, err := jp.Search(r.metadata)
	if err != nil {
		log.Panic(err)
	}
	return v
}

func (r *Rec) GetString(path string) string {
	val := r.Get(path)
	if str, ok := val.(string); ok {
		return str
	}
	return ""
}

func (r *Rec) GetStringSlice(path string) (strs []string) {
	vals, ok := r.Get(path).([]interface{})
	if !ok {
		return
	}
	for _, v := range vals {
		if str, ok := v.(string); ok {
			strs = append(strs, str)
		}
	}
	return
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

func (e *engine) AddRepresentation(recID, name string, data []byte) error {
	return e.store.AddRepresentation(recID, name, data)
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
