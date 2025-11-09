package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/models/remove"
)

func (t Transpiler) TranspileRemove(r remove.Remove) (error, string) {
	query := fmt.Sprintf("DELETE FROM \n%s", formatter.Indent(r.Target))

	conditionString := ""
	if len(r.Conditions) > 0 {
		var conditionError error
		conditionError, conditionString = t.TranspileConditions(r.Conditions, false)
		if conditionError != nil {
			return conditionError, ""
		}
	}

	if conditionString != "" {
		query += conditionString
	}

	return nil, query + ";"
}
