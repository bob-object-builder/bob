package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/insert"
)

func TranspileInsert(n *insert.Insert) (*failure.Failure, string) {
	selected := formatter.Indent(
		fmt.Sprintf("(%s)", n.Fields.Join(", ")),
	)

	values := value.NewArray[string]()

	for _, v := range n.Values {
		row := v.Join(", ")

		if n.IsBulk {
			if v.Length() < n.Fields.Length() {
				return failure.NotEnoughValues(n.Target), ""
			}

			values.Push(
				formatter.Indent(fmt.Sprintf("(%s)", row)),
			)
		} else {
			values.Push(
				formatter.Indent(row),
			)
		}
	}

	var valuesString string
	var sql string

	if n.IsBulk {
		valuesString = values.Join(",\n")
		sql = "INSERT INTO %s\n%s\nVALUES\n%s;"
	} else {
		valuesString = values.Join(",\n")
		sql = "INSERT INTO %s\n%s\nVALUES (\n%s\n);"
	}

	return nil, fmt.Sprintf(sql, n.Target, selected, valuesString)
}
