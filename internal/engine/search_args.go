package engine

type Filters map[string][]string

type SearchArgs struct {
	Query   string  `form:"q"`
	Filters Filters `form:"f"`
	Size    int     `form:"size"`
	Skip    int     `form:"skip"`
}

func (a SearchArgs) WithFilter(field string, terms ...string) SearchArgs {
	if a.Filters == nil {
		a.Filters = make(Filters)
	}

	a.Filters[field] = terms

	return a
}
