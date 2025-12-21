package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/model/remove"
)

func TranspileDelete(v *remove.Delete) (*failure.Failure, string) {
	sql := fmt.Sprintf("DELETE FROM %s", v.Target)

	conditionsError, conditions := TranspileCondition(&v.Conditions, false)
	if conditionsError != nil {
		return conditionsError, ""
	}

	sql += conditions
	sql += ";"

	return nil, sql
}
