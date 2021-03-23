package engine

type Scope map[string][]string

type SearchArgs struct {
	Scope Scope  `form:"f"`
	Query string `form:"q"`
	Size  int    `form:"size"`
	Skip  int    `form:"skip"`
}

func (a SearchArgs) WithScope(scope Scope) SearchArgs {
	if a.Scope == nil {
		a.Scope = make(Scope)
	}

	for field, terms := range scope {
		a.Scope[field] = terms
	}

	return a
}
