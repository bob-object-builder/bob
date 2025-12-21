package join

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/value"
)

type JoinSQL struct {
	Table string
}

func (jc *Joins) ToSQL(parenTable string) (*failure.Failure, string, value.Array[Selection]) {
	joins := value.NewArray[string]()
	selections := value.NewArray[Selection]()

	for _, join := range jc.Joins {
		table := join.ReferencedTable
		column := join.ReferencedColumn

		parent := fmt.Sprintf("%s.%s_%s", table, parenTable, column)

		query := fmt.Sprintf(
			"LEFT JOIN %s ON %s = %s.%s",
			table,
			parent,
			parenTable,
			column,
		)

		selections.Push(join.Selected...)
		joins.Push(query)
	}

	return nil, joins.Join("\n"), *selections
}
