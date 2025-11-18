package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/insert"
	"strings"
)

func (t Transpiler) TranspileInsert(i insert.Insert) string {
	query := "INSERT INTO\n"
	query += formatter.Indent("%s(%s)\nVALUES\n%s;")
	columns := strings.Join(i.Columns, ", ")
	accumulatedValues := array.New[string]()

	for _, values := range i.Values {
		accumulatedValues.Push(
			fmt.Sprintf(
				formatter.Indent("(%s)"),
				strings.Join(values, ", "),
			),
		)
	}

	return fmt.Sprintf(query, i.Target, columns, strings.Join(*accumulatedValues, ",\n"))
}
