package action

import (
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/models/action/get"
	"salvadorsru/bob/internal/models/action/insert"
	"salvadorsru/bob/internal/models/action/remove"
	"salvadorsru/bob/internal/models/action/set"
)

func Parse(actions lexer.Blocks) []any {
	queries := []any{}

	for _, action := range actions {
		switch action.Command {
		case lexer.Get:
			query := get.NewQuery(action)
			queries = append(queries, query)
		case lexer.New:
			query := insert.NewQuery(action)
			queries = append(queries, query)
		case lexer.Set:
			query := set.NewQuery(action)
			queries = append(queries, query)
		case lexer.Delete:
			query := remove.NewQuery(action)
			queries = append(queries, query)
		}
	}

	return queries
}
