package records

import (
	"encoding/json"
	"time"
)

type Rec struct {
	ID          string          `json:"id"`
	Type        string          `json:"type"`
	Collection  []string        `json:"collection"`
	Title       string          `json:"title"`
	RawMetadata json.RawMessage `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	RawSource   json.RawMessage `json:"source"`
}

type Scope map[string][]string

type SearchArgs struct {
	Scope Scope  `form:"-"`
	Query string `form:"q"`
	Size  int    `form:"size"`
	Skip  int    `form:"skip"`
}

type Hit struct {
	Rec
	RawHighlight json.RawMessage `json:"highlight"`
}

type Hits struct {
	Total int    `json:"total"`
	Hits  []*Hit `json:"hits"`
}

func (s SearchArgs) WithScope(scope Scope) SearchArgs {
	// TODO merge scope new scope with old scope
	s.Scope = scope
	return s
}
