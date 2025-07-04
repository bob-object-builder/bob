package table

import (
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func Transpile(driver drivers.Driver, tables lexer.Tables) (*utils.Object[*Table], string) {
	parsedTables := Parse(tables)
	queries := []string{}

	for _, tableName := range parsedTables.Order {
		table := parsedTables.Get(tableName)
		queries = append(queries, table.ToQuery(driver))
	}

	if len(queries) == 0 {
		return nil, ""
	}

	return &parsedTables, strings.Join(queries, ";\n\n") + ";"
}
