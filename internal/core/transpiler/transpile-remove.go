package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/models/remove"
)

func (t Transpiler) TranspileRemove(r remove.Remove) string {
	query := fmt.Sprintf("DELETE FROM \n%s", formatter.Indent(r.Target))

	conditionString := ""
	if len(r.Conditions) > 0 {
		conditionString = "\nWHERE\n" + t.TranspileConditions(r.Conditions)
	}

	if conditionString != "" {
		query += conditionString
	}

	return query + ";"
}
