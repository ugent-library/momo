package momo

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v6"
)

type Recs struct {
	es    *elasticsearch.Client
	index string
}

func (r *Recs) AutocompleteSearch(qs string) map[string]interface{} {
	var buf bytes.Buffer
	var query map[string]interface{}

	if len(qs) == 0 {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	} else {
		query = map[string]interface{}{
			"query": map[string]interface{}{
				"match": map[string]interface{}{
					"ac": qs,
				},
			},
		}
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}
	res, err := r.es.Search(
		r.es.Search.WithContext(context.Background()),
		r.es.Search.WithIndex(r.index),
		r.es.Search.WithBody(&buf),
		r.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	var resData map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	log.Printf("es response: %s", resData)
	return resData["hits"].(map[string]interface{})
}
