package records

import (
	"encoding/json"
	"time"
)

type Rec struct {
	ID         string          `json:"id"`
	Type       string          `json:"type"`
	Collection []string        `json:"collection"`
	Title      string          `json:"title"`
	Metadata   json.RawMessage `json:"metadata"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
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
	Highlight json.RawMessage `json:"highlight"`
}

type Hits struct {
	Total int    `json:"total"`
	Hits  []*Hit `json:"hits"`
}

// TODO merge scope new scope with old scope
func (s SearchArgs) WithScope(scope Scope) SearchArgs {
	s.Scope = scope
	return s
}
