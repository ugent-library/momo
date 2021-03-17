package engine

type Scope map[string][]string

type SearchArgs struct {
	Scope Scope  `form:"-"`
	Query string `form:"q"`
	Size  int    `form:"size"`
	Skip  int    `form:"skip"`
}

func (s SearchArgs) WithScope(scope Scope) SearchArgs {
	// TODO merge new scope with old scope
	s.Scope = scope
	return s
}
