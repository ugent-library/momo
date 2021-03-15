package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/ugent-library/momo/engine"
)

type Store struct {
	Client       *elasticsearch.Client
	IndexName    string
	IndexMapping string
}

func (s *Store) CreateRecIndex() error {
	r := strings.NewReader(s.IndexMapping)
	res, err := s.Client.Indices.Create(s.IndexName, s.Client.Indices.Create.WithBody(r))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (s *Store) DeleteRecIndex() error {
	res, err := s.Client.Indices.Delete([]string{s.IndexName})
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (s *Store) AddRec(rec *engine.Rec) error {
	payload, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index:      s.IndexName,
		DocumentID: rec.ID,
		Body:       bytes.NewReader(payload),
	}.Do(ctx, s.Client)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return err
		}
		return fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	return nil
}

// TODO don't die
func (s *Store) AddRecs(c <-chan *engine.Rec) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  s.IndexName,
		Client: s.Client,
		// NumWorkers: 4, // The number of worker goroutines (default NumCPU)
		// FlushBytes: int(flushBytes), // The flush threshold in bytes (default 5MB)
		// FlushInterval: 30 * time.Second, // The periodic flush interval (default 30S)
		OnError: func(c context.Context, e error) {
			log.Printf("ERROR: %s", e)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// var countSuccessful uint64

	// start := time.Now().UTC()

	for rec := range c {

		payload, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
		}

		// Add item to the BulkIndexer
		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:       "index",
				DocumentID:   rec.ID,
				DocumentType: "_doc",
				Body:         bytes.NewReader(payload),
				// OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				// atomic.AddUint64(&countSuccessful, 1)
				// },
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Fatalf("Unexpected error: %s", err)
		}

	}

	// Close the indexer
	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}

	// biStats := bi.Stats()

	// dur := time.Since(start)

	// if biStats.NumFailed > 0 {
	// 	log.Fatalf(
	// 		"Indexed [%d] documents with [%d] errors in %s (%d docs/sec)",
	// 		int64(biStats.NumFlushed),
	// 		int64(biStats.NumFailed),
	// 		dur.Truncate(time.Millisecond),
	// 		int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
	// 	)
	// } else {
	// 	log.Printf(
	// 		"Sucessfuly indexed [%d] documents in %s (%d docs/sec)",
	// 		int64(biStats.NumFlushed),
	// 		dur.Truncate(time.Millisecond),
	// 		int64(1000.0/float64(dur/time.Millisecond)*float64(biStats.NumFlushed)),
	// 	)
	// }
}

func (s *Store) SearchRecs(args engine.SearchArgs) (*engine.RecHits, error) {
	var buf bytes.Buffer
	var query map[string]interface{}

	if len(args.Query) == 0 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"multi_match": map[string]interface{}{
					"query":    args.Query,
					"fields":   []string{"id^100", "metadata.title.ngram"},
					"operator": "and",
				},
			},
		}
	}

	if args.Scope != nil {
		var terms []map[string]interface{}

		for k, v := range args.Scope {
			terms = append(terms, map[string]interface{}{"terms": map[string]interface{}{k: v}})
		}

		query = map[string]interface{}{
			"query": map[string]interface{}{
				"bool": map[string]interface{}{
					"must": query["query"],
					"filter": map[string]interface{}{
						"bool": map[string]interface{}{
							"must": terms,
						},
					},
				},
			},
		}
	}

	query["size"] = args.Size
	query["from"] = args.Skip
	query["highlight"] = map[string]interface{}{
		"require_field_match": false,
		"pre_tags":            []string{"<mark>"},
		"post_tags":           []string{"</mark>"},
		"fields": map[string]interface{}{
			"metadata.title.ngram": map[string]interface{}{},
		},
	}
	query["aggs"] = map[string]interface{}{
		"collection": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "collection",
			},
		},
		"type": map[string]interface{}{
			"terms": map[string]interface{}{
				"field": "type",
			},
		},
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := s.Client.Search(
		s.Client.Search.WithContext(context.Background()),
		s.Client.Search.WithIndex(s.IndexName),
		s.Client.Search.WithBody(&buf),
		s.Client.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	type resEnvelope struct {
		Hits struct {
			Total int
			Hits  []struct {
				Source    json.RawMessage `json:"_source"`
				Highlight json.RawMessage
			}
		}
		Aggregations json.RawMessage
	}

	var r resEnvelope
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	hits := engine.RecHits{
		Total: r.Hits.Total,
		Hits:  []*engine.RecHit{},
	}

	if len(r.Aggregations) > 0 {
		hits.RawAggregation = r.Aggregations
	}

	for _, h := range r.Hits.Hits {
		var hit engine.RecHit

		if err := json.Unmarshal(h.Source, &hit); err != nil {
			return nil, err
		}

		if len(h.Highlight) > 0 {
			hit.RawHighlight = h.Highlight
		}

		hits.Hits = append(hits.Hits, &hit)
	}

	return &hits, nil
}
