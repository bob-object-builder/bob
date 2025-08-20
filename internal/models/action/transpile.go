package action

import (
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/models/action/get"
	"salvadorsru/bob/internal/models/action/insert"
	"salvadorsru/bob/internal/models/action/set"
	"strings"
)

func Transpile(driver drivers.Driver, blocks lexer.Blocks) (error, string) {
	queries := []string{}
	actions := Parse(blocks)

	for _, action := range actions {
		switch v := action.(type) {
		case get.Get:
			queryError, query := v.ToQuery(driver)
			if queryError != nil {
				return queryError, ""
			}

			queries = append(queries, query)
		case insert.New:
			query := v.ToQuery(driver)
			queries = append(queries, query)
		case set.Set:
			query := v.ToQuery(driver)
			queries = append(queries, query)
		}
	}

	if len(queries) == 0 {
		return nil, ""
	}

	return nil, strings.Join(queries, ";\n\n") + ";"
}
