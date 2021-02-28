package listing

type SearchScope map[string][]string

type SearchArgs struct {
	Scope SearchScope `form:"-"`
	Query string      `form:"q"`
	Size  int         `form:"size"`
	Skip  int         `form:"skip"`
}

// TODO merge scope new scope with old scope
func (s SearchArgs) WithScope(scope SearchScope) SearchArgs {
	s.Scope = scope
	return s
}
