package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/models/remove"
)

func (t Transpiler) TranspileRemove(r remove.Remove) (*failure.Failure, string) {
	query := fmt.Sprintf("DELETE FROM \n%s", formatter.Indent(r.Target))

	conditionString := ""
	if len(r.Conditions) > 0 {
		var conditionError *failure.Failure
		conditionError, conditionString = t.TranspileConditions(r.Conditions, false)
		if conditionError != nil {
			return conditionError, ""
		}
	} else {
		if !r.RemoveAll {
			return failure.DeleteCondition(r.Target), ""
		}
	}

	if conditionString != "" {
		query += conditionString
	}

	return nil, query + ";"
}
