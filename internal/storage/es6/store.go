package es6

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/elastic/go-elasticsearch/v6/esapi"
	"github.com/elastic/go-elasticsearch/v6/esutil"
	"github.com/pkg/errors"
	"github.com/ugent-library/momo/internal/engine"
)

type Config struct {
	ClientConfig elasticsearch.Config
	IndexPrefix  string
	IndexMapping string
}

type store struct {
	client *elasticsearch.Client
	Config
}

type esRec struct {
	ID         string                 `json:"id"`
	Collection string                 `json:"collection"`
	Type       string                 `json:"type"`
	Metadata   map[string]interface{} `json:"metadata"`
	Source     string                 `json:"source,omitempty"`
	SourceID   string                 `json:"source_id,omitempty"`
	CreatedAt  string                 `json:"createdAt"`
	UpdatedAt  string                 `json:"updatedAt"`
}

type resEnvelope struct {
	ScrollID string `json:"_scroll_id"`
	Hits     struct {
		Total int
		Hits  []struct {
			Source    json.RawMessage `json:"_source"`
			Highlight json.RawMessage
		}
	}
	Aggregations json.RawMessage
}

type M map[string]interface{}

func New(c Config) (*store, error) {
	client, err := elasticsearch.NewClient(c.ClientConfig)
	if err != nil {
		return nil, err
	}
	return &store{client, c}, nil
}

func (s *store) indexName(idx string) string {
	return s.IndexPrefix + idx
}

func (s *store) CreateRecIndex() error {
	r := strings.NewReader(s.IndexMapping)
	res, err := s.client.Indices.Create(s.indexName("rec"), s.client.Indices.Create.WithBody(r))
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (s *store) DeleteRecIndex() error {
	res, err := s.client.Indices.Delete([]string{s.indexName("rec")})
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (s *store) Reset() error {
	res, err := s.client.Indices.Exists([]string{s.indexName("rec")})
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200:
		err = s.DeleteRecIndex()
		if err != nil {
			return err
		}
	case 404:
		// index doesn't exist, do nothing
	default:
		if res.IsError() {
			return fmt.Errorf("error: %s", res)
		}
	}
	return s.CreateRecIndex()
}

func (s *store) AddRec(rec *engine.Rec) error {
	payload, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	ctx := context.Background()
	res, err := esapi.CreateRequest{
		Index:      s.indexName("rec"),
		DocumentID: rec.ID,
		Body:       bytes.NewReader(payload),
	}.Do(ctx, s.client)
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

func (s *store) AddRecs(c <-chan *engine.Rec) {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  s.indexName("rec"),
		Client: s.client,
		OnError: func(c context.Context, e error) {
			log.Printf("ERROR: %s", e)
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for rec := range c {
		// not needed anymore in es7 with date nano type
		r := esRec{
			ID:         rec.ID,
			Collection: rec.Collection,
			Type:       rec.Type,
			Metadata:   rec.Metadata,
			Source:     rec.Source,
			SourceID:   rec.SourceID,
			CreatedAt:  rec.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:  rec.UpdatedAt.UTC().Format(time.RFC3339),
		}

		payload, err := json.Marshal(&r)
		if err != nil {
			log.Panic(err)
		}

		err = bi.Add(
			context.Background(),
			esutil.BulkIndexerItem{
				Action:       "index",
				DocumentID:   rec.ID,
				DocumentType: "_doc",
				Body:         bytes.NewReader(payload),
				OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
					if err != nil {
						log.Panicf("ERROR: %s", err)
					} else {
						log.Panicf("ERROR: %s: %s", res.Error.Type, res.Error.Reason)
					}
				},
			},
		)

		if err != nil {
			log.Panicf("Unexpected error: %s", err)
		}
	}

	// Close the indexer
	if err := bi.Close(context.Background()); err != nil {
		log.Panicf("Unexpected error: %s", err)
	}
}

func (s *store) SearchRecs(args engine.SearchArgs) (*engine.RecHits, error) {
	query, queryFilter, termsFilters := buildQuery(args)

	query["size"] = args.Size
	query["from"] = args.Skip

	if args.Highlight {
		query["highlight"] = M{
			"require_field_match": false,
			"pre_tags":            []string{"<mark>"},
			"post_tags":           []string{"</mark>"},
			"fields": M{
				"metadata.title.ngram":       M{},
				"metadata.author.name.ngram": M{},
			},
		}
	}

	if len(args.Facets) > 0 {
		query["aggs"] = M{
			"facets": M{
				"global": M{},
				"aggs":   M{},
			},
		}

		// facet filter contains all query and all filters except itself
		for _, field := range args.Facets {
			filters := []M{queryFilter}

			for _, filter := range termsFilters {
				if _, found := filter["terms"].(M)[field]; found {
					continue
				} else {
					filters = append(filters, filter)
				}
			}

			query["aggs"].(M)["facets"].(M)["aggs"].(M)[field] = M{
				"filter": M{"bool": M{"must": filters}},
				"aggs": M{
					"facet": M{
						"terms": M{
							"field":         field,
							"min_doc_count": 1,
							"order":         M{"_key": "asc"},
							"size":          200,
						},
					},
				},
			}
		}
	}

	opts := []func(*esapi.SearchRequest){
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.indexName("rec")),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithSort("_doc"),
	}

	if args.Cursor {
		opts = append(opts, s.client.Search.WithScroll(time.Minute))
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	opts = append(opts, s.client.Search.WithBody(&buf))

	res, err := s.client.Search(opts...)
	if err != nil {
		return nil, err
	}

	return decodeRes(res)
}

func (s *store) SearchMoreRecs(cursorID string) (*engine.RecHits, error) {
	// use search after api in es7
	res, err := s.client.Scroll(
		s.client.Scroll.WithScrollID(cursorID),
		s.client.Scroll.WithScroll(time.Minute),
	)

	if err != nil {
		return nil, err
	}

	return decodeRes(res)
}

func decodeRes(res *esapi.Response) (*engine.RecHits, error) {
	defer res.Body.Close()

	if res.IsError() {
		buf := new(strings.Builder)
		if _, err := io.Copy(buf, res.Body); err != nil {
			return nil, err
		}
		return nil, errors.New("Es6 error response: " + buf.String())
	}

	var r resEnvelope
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, errors.Wrap(err, "Error parsing the response body")
	}

	hits := engine.RecHits{
		CursorID: r.ScrollID,
		Total:    r.Hits.Total,
		Hits:     []*engine.RecHit{},
	}

	if len(r.Aggregations) > 0 {
		hits.RawAggregation = r.Aggregations
	}

	for _, h := range r.Hits.Hits {
		var hit engine.RecHit

		if err := json.Unmarshal(h.Source, &hit.Rec); err != nil {
			return nil, err
		}

		if len(h.Highlight) > 0 {
			hit.RawHighlight = h.Highlight
		}

		hits.Hits = append(hits.Hits, &hit)
	}

	return &hits, nil
}
