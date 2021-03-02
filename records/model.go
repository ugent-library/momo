package records

import (
	"encoding/json"
	"time"
)

// TODO make a struct for new recs without timestamps?
type Rec struct {
	ID         string   `json:"id"`
	Type       string   `json:"type"`
	Collection []string `json:"collection"`
	Title      string   `json:"title"`
	Metadata   json.RawMessage
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Scope map[string][]string

type SearchArgs struct {
	Scope Scope  `form:"-"`
	Query string `form:"q"`
	Size  int    `form:"size"`
	Skip  int    `form:"skip"`
}

type Hits struct {
	Total int    `json:"total"`
	Hits  []*Rec `json:"hits"`
}

// TODO merge scope new scope with old scope
func (s SearchArgs) WithScope(scope Scope) SearchArgs {
	s.Scope = scope
	return s
}
