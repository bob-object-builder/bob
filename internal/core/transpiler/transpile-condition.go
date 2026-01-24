package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
	"strings"
)

func TranspileCondition(cs *condition.Conditions, isGrouped bool, aliases *value.Map[string]) (*failure.Failure, string) {
	if len(cs.Conditions) == 0 {
		return nil, ""
	}

	var conditions = value.NewArray[string]()

	for i, c := range cs.Conditions {
		var conditionKey string
		if i > 0 {
			if c.Condition == "if" {
				conditionKey = "AND"
			} else {
				conditionKey = "OR"
			}
		}

		if c.Comparator == "" {
			return failure.MalformedCondition(c.Target), ""
		}

		if c.And.Length() == 0 && c.Else.Length() == 0 {
			return failure.ConditionValidation(c.From, c.Target), ""
		}

		var leftSide string
		hasDot := strings.Contains(c.Target, ".")
		useAlias := aliases != nil && aliases.Get(c.Target) == nil

		switch {
		case hasDot:
			leftSide = c.Target
		case useAlias:
			leftSide = fmt.Sprintf("%s.%s", c.From, c.Target)
		default:
			leftSide = c.Target
		}

		operation := fmt.Sprintf("%s %s", leftSide, c.Comparator)

		and := c.And.Join(fmt.Sprintf("\nAND %s ", operation))
		and = fmt.Sprintf("%s %s", operation, and)

		var or string
		if len(c.Else) > 0 {
			or = c.Else.Join(fmt.Sprintf("\nOR %s ", operation))
			or = fmt.Sprintf("\nOR %s %s", operation, or)
		}

		full := fmt.Sprintf("%s %s %s", conditionKey, and, or)
		conditions.Push(strings.TrimSpace(full))
	}

	head := "\nWHERE\n"
	if isGrouped {
		head = "HAVING\n"
	}

	return nil, head + formatter.IndentLines(conditions.Join("\n"))
}
