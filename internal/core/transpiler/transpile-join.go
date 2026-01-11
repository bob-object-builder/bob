package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
	"salvadorsru/bob/internal/model/join"
)

func TranspileJoin(j *join.Joins, parenTable string) (*failure.Failure, string, value.Array[join.Selection], value.Array[condition.Condition]) {
	joins := value.NewArray[string]()
	selections := value.NewArray[join.Selection]()
	conditions := value.NewArray[condition.Condition]()

	for _, join := range j.Joins {
		referencedTable := join.ReferencedTable
		column := join.ReferencedColumn

		query := fmt.Sprintf(
			"LEFT JOIN %s ON %s.%s_%s = %s.%s",
			referencedTable,
			parenTable,
			referencedTable,
			column,
			referencedTable,
			column,
		)

		selections.Push(join.Selected...)
		joins.Push(query)
		conditions.Push(join.Conditions.Conditions...)
	}

	return nil, joins.Join("\n"), *selections, *conditions
}
