package listing

type SearchScope map[string][]string

type SearchArgs struct {
	Query   string
	Scope   SearchScope
	Page    int
	PerPage int
}

func (s SearchArgs) WithScope(scope SearchScope) SearchArgs {
	s.Scope = scope
	return s
}
