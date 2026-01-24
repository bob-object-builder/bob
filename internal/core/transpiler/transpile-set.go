package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/set"
)

func TranspileSet(s *set.Set) (*failure.Failure, string) {
	sql := "UPDATE %s\nSET\n%s%s;"
	fields := value.NewArray[string]()

	for v := range s.Values.Range() {
		fields.Push(formatter.Indent(fmt.Sprintf("%s = %s", v.Key, v.Value)))
	}

	conditionsError, conditions := TranspileCondition(&s.Conditions, false, nil)
	if conditionsError != nil {
		return conditionsError, ""
	}

	return nil, fmt.Sprintf(sql, formatter.ToReferenceCase(s.Target), fields.Join(",\n"), conditions)
}
