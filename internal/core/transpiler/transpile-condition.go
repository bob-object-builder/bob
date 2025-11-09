package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/lib/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/condition"
	"strings"
)

func (t Transpiler) TranspileConditions(conds array.Array[condition.Condition], isGrouped bool) (*failure.Failure, string) {
	if len(conds) == 0 {
		return nil, ""
	}

	var conditions = array.New[string]()

	for i, c := range conds {
		var conditionKey string
		if i > 0 {
			if c.Condition == condition.If {
				conditionKey = "AND"
			} else {
				conditionKey = "OR"
			}
		}

		if c.Comparator == "" {
			return failure.MalformedCondition(c.Target), ""
		}

		if c.And.Length() == 0 && c.Else.Length() == 0 {
			return failure.ConditionValidation(c.Table, c.Target), ""
		}

		var operation string
		if strings.Contains(c.Target, ".") {
			operation = fmt.Sprintf("%s %s", c.Target, c.Comparator)
		} else {
			operation = fmt.Sprintf("%s.%s %s", c.Table, c.Target, c.Comparator)
		}

		and := strings.Join(c.And, fmt.Sprintf("\nAND %s ", operation))
		and = fmt.Sprintf("%s %s", operation, and)

		var or string
		if len(c.Else) > 0 {
			or = strings.Join(c.Else, fmt.Sprintf("\nOR %s ", operation))
			or = fmt.Sprintf("\nOR %s %s", operation, or)
		}

		full := fmt.Sprintf("%s %s %s", conditionKey, and, or)
		conditions.Push(strings.TrimSpace(full))
	}

	head := "\nWHERE\n"
	if isGrouped {
		head = "\nHAVING\n"
	}

	return nil, head + formatter.IndentLines(strings.Join(*conditions, "\n"))
}
