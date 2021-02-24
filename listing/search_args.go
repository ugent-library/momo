package listing

type SearchScope map[string][]string

type SearchArgs struct {
	Query string
	Scope SearchScope
}

func (s SearchArgs) WithScope(scope SearchScope) SearchArgs {
	s.Scope = scope
	return s
}
