package es6

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
	"github.com/ugent-library/momo/listing"
)

type Store struct {
	Client       *elasticsearch.Client
	IndexName    string
	IndexMapping string
}

func (s *Store) CreateIndex() error {
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

func (s *Store) DeleteIndex() error {
	res, err := s.Client.Indices.Delete([]string{s.IndexName})
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (s *Store) SearchRecs(args listing.SearchArgs) (*listing.RecHits, error) {
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
				"match": map[string]interface{}{
					"ac": args.Query,
				},
			},
		}
	}

	// clear as mud
	if args.Scope != nil {
		terms := make([]map[string]interface{}, len(args.Scope))
		for k, v := range args.Scope {
			terms = append(terms, map[string]interface{}{"terms": map[string]interface{}{k: v}})
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
	}

	query["size"] = args.Size
	query["from"] = args.Skip

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := s.Client.Search(
		// s.Es.Search.WithContext(context.Background()),
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
		Took int
		Hits struct {
			Total int
			Hits  []struct {
				ID     string          `json:"_id"`
				Source json.RawMessage `json:"_source"`
				// Highlights json.RawMessage `json:"highlight"`
				// Sort []interface{} `json:"sort"`
			}
		}
	}

	type resHit struct {
		Title string `json:"title"`
	}

	var r resEnvelope
	if err = json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	hits := listing.RecHits{}
	hits.Total = r.Hits.Total

	if len(r.Hits.Hits) == 0 {
		hits.Hits = []*listing.RecHit{}
		return &hits, nil
	}

	for _, hit := range r.Hits.Hits {
		var rh resHit
		var h listing.RecHit
		h.ID = hit.ID

		if err := json.Unmarshal(hit.Source, &rh); err != nil {
			return nil, err
		}

		h.Title = rh.Title

		// if len(hit.Highlights) > 0 {
		// 	if err := json.Unmarshal(hit.Highlights, &h.Highlights); err != nil {
		// 		return &results, err
		// 	}
		// }

		hits.Hits = append(hits.Hits, &h)
	}

	return &hits, nil
}