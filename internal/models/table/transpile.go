package table

import (
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"strings"
)

func Transpile(driver drivers.Driver, tables lexer.Tables) (error, string) {
	parsedTablesError, parsedTables := Parse(tables)
	if parsedTablesError != nil {
		return parsedTablesError, ""
	}
	queries := []string{}

	for _, tableName := range parsedTables.Order {
		table := parsedTables.Get(tableName)

		err, query := table.ToQuery(driver)
		if err != nil {
			return err, ""
		}

		queries = append(queries, query)
	}

	if len(queries) == 0 {
		return nil, ""
	}

	return nil, strings.Join(queries, ";\n\n") + ";"
}
