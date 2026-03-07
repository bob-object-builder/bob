package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
	"strings"
)

func TranspileChecks(cs *condition.Conditions, isGrouped bool, aliases *value.Map[string]) (*failure.Failure, string) {
	if len(cs.Conditions) == 0 {
		return nil, ""
	}

	var conditions = value.NewArray[string]()

	for _, c := range cs.Conditions {

		if c.And.Length() == 0 && c.Else.Length() == 0 {
			return failure.ConditionValidation(c.From, c.Target), ""
		}

		leftSide := c.Target

		operation := fmt.Sprintf("%s %s", leftSide, c.Comparator)

		and := c.And.Join(fmt.Sprintf(" AND %s ", operation))
		and = fmt.Sprintf("%s %s", operation, and)

		var or string
		if len(c.Else) > 0 {
			or = c.Else.Join(fmt.Sprintf(" OR %s ", operation))
			or = fmt.Sprintf("OR %s %s", operation, or)
		}

		full := strings.TrimSpace(fmt.Sprintf("%s %s", and, or))
		full = fmt.Sprintf("CHECK(%s)", full)
		conditions.Push(strings.TrimSpace(full))
	}

	return nil, "\n" + formatter.IndentLines(conditions.Join(",\n"))
}
