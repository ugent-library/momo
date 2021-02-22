package store

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v6"
)

type Es struct {
	Client       *elasticsearch.Client
	IndexName    string
	IndexMapping string
}

func (s *Es) CreateIndex() error {
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

func (s *Es) DeleteIndex() error {
	res, err := s.Client.Indices.Delete([]string{s.IndexName})
	if err != nil {
		return err
	}
	if res.IsError() {
		return fmt.Errorf("error: %s", res)
	}
	return nil
}

func (s *Es) SearchRecs(qs string) (map[string]interface{}, error) {
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

	var resData map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&resData); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	// log.Printf("es response: %s", resData)
	return resData["hits"].(map[string]interface{}), nil
}
