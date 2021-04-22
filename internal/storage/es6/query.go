package es6

import (
	"github.com/ugent-library/momo/internal/engine"
)

func buildQuery(args engine.SearchArgs) (M, M, []M) {
	var query M
	var queryFilter M
	var termsFilters []M

	if len(args.Query) == 0 {
		queryFilter = M{
			"match_all": M{},
		}
	} else {
		queryFilter = M{
			"multi_match": M{
				"query":    args.Query,
				"fields":   []string{"id^100", "metadata.identifier^50", "metadata.title.ngram", "metadata.author.name.ngram"},
				"operator": "and",
			},
		}
	}

	if args.Filters == nil {
		query = M{"query": queryFilter}
	} else {
		for field, terms := range args.Filters {
			termsFilters = append(termsFilters, M{"terms": M{field: terms}})
		}

		query = M{
			"query": M{
				"bool": M{
					"must": queryFilter,
					"filter": M{
						"bool": M{
							"must": termsFilters,
						},
					},
				},
			},
		}
	}

	query["size"] = args.Size
	query["from"] = args.Skip

	return query, queryFilter, termsFilters
}
