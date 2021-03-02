package listing

import (
	"encoding/json"
	"time"
)

type RecHits struct {
	Total int       `json:"total"`
	Hits  []*RecHit `json:"hits"`
}

type RecHit struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type Rec struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Collection []string
	Title      string `json:"title"`
	Metadata   json.RawMessage
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
