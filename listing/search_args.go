package listing

type SearchScope map[string][]string

type SearchArgs struct {
	Scope SearchScope `schema:"-"`
	Query string      `schema:"q"`
	Size  int         `schema:"size"`
	Skip  int         `schema:"skip"`
}

// TODO merge scope new scope with old scope
func (s SearchArgs) WithScope(scope SearchScope) SearchArgs {
	s.Scope = scope
	return s
}
