package listing

type Scope map[string][]string

type SearchArgs struct {
	Scope Scope  `form:"-"`
	Query string `form:"q"`
	Size  int    `form:"size"`
	Skip  int    `form:"skip"`
}

// TODO merge scope new scope with old scope
func (s SearchArgs) WithScope(scope Scope) SearchArgs {
	s.Scope = scope
	return s
}
