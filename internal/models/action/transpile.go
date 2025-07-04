package action

import (
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"salvadorsru/bob/internal/models/action/get"
	"salvadorsru/bob/internal/models/action/insert"
	"salvadorsru/bob/internal/models/table"
	"strings"
)

func Transpile(driver drivers.Driver, tables *utils.Object[*table.Table], blocks lexer.Blocks) string {
	queries := []string{}
	actions := Parse(blocks)

	for _, action := range actions {
		switch v := action.(type) {
		case get.Get:
			query := v.ToQuery(driver)
			queries = append(queries, query)
		case insert.New:
			query := v.ToQuery(driver)
			queries = append(queries, query)
		}
	}

	if len(queries) == 0 {
		return ""
	}

	return strings.Join(queries, ";\n\n") + ";"
}
