package engine

import "strings"

type Scope map[string][]string

type SearchArgs struct {
	Scope      Scope  `form:"-"`
	Query      string `form:"q"`
	Size       int    `form:"size"`
	Skip       int    `form:"skip"`
	Type       string `form:"type"`
	Collection string `form:"collection"`
}

func (s SearchArgs) WithScope(scope Scope) SearchArgs {
	// TODO merge new scope with old scope
	s.Scope = scope

	if len(s.Type) > 0 {
		s.Scope["type"] = strings.Split(s.Type, "-")
	} else {
		delete(s.Scope, "type")
	}

	if len(s.Collection) > 0 {
		s.Scope["collection"] = strings.Split(s.Collection, "-")
	} else {
		delete(s.Scope, "collection")
	}

	return s
}
