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

// TODO don't die
// TODO send errors back over a channel
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
		payload, err := json.Marshal(rec)
		if err != nil {
			log.Fatal(err)
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
}

func (s *store) SearchRecs(args engine.SearchArgs) (*engine.RecHits, error) {
	facetFields := []string{"collection", "type"}
	var buf bytes.Buffer
	var query M
	var queryFilter M
	var termsFilters []M

	if len(args.Query) == 0 {
		queryFilter = M{
			"match_all": M{},
		}
	} else {
		queryFilter = M{
			"multi_match": M{
				"query":    args.Query,
				"fields":   []string{"id^100", "metadata.title.ngram"},
				"operator": "and",
			},
		}
	}

	if args.Scope == nil {
		query = M{"query": queryFilter}
	} else {
		for field, terms := range args.Scope {
			termsFilters = append(termsFilters, M{"terms": M{field: terms}})
		}

		query = M{
			"query": M{
				"bool": M{
					"must": queryFilter,
					"filter": M{
						"bool": M{
							"must": termsFilters,
						},
					},
				},
			},
		}
	}

	query["size"] = args.Size
	query["from"] = args.Skip
	query["highlight"] = M{
		"require_field_match": false,
		"pre_tags":            []string{"<mark>"},
		"post_tags":           []string{"</mark>"},
		"fields": M{
			"metadata.title.ngram": M{},
		},
	}
	query["aggs"] = M{
		"facets": M{
			"global": M{},
			"aggs":   M{},
		},
	}

	// facet filter contains all query and all filters except itself
	for _, field := range facetFields {
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
						"size":          500, // TODO give better value or use nested facets or composite aggregation
					},
				},
			},
		}
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.indexName("rec")),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
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
