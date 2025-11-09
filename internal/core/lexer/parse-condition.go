package lexer

import (
	"fmt"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/function"
	"strings"
)

func (l *Lexer) ParseCondition(target string) condition.Condition {
	tokens := l.tokens

	errorCondition := condition.Condition{
		Table: target,
	}

	if len(tokens) == 0 {
		return errorCondition
	}

	conditionKey, tokens := tokens[0], tokens[1:]

	newCondition := &condition.Condition{
		Table:     target,
		Condition: condition.ConditionKey(conditionKey),
	}

	var (
		currentList   = &newCondition.And
		parsingString bool
		buffer        = array.New[string]()
	)

	prefix := func(str string) string {
		return function.PrefixParameters(fmt.Sprintf("%s.", target), str)
	}

	push := func(str string) {
		if newCondition.Target == "" {
			newCondition.Target = prefix(str)
		} else {
			currentList.Push(prefix(str))
		}
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch {
		case parsingString || checker.IsStringStart(token):
			parsingString = true
			buffer.Push(token)
			if checker.IsStringEnd(token) {
				parsingString = false
				fullString := strings.Join(*buffer, " ")
				buffer.Clean()
				push(formatter.NormalizeString(fullString))
			}

		case function.IsFunctionStart(token):
			completeFunction, lastIndex := function.ReconstructFunction(tokens, i)
			i = lastIndex
			push(completeFunction)

		case condition.IsAnd(token):
			currentList = &newCondition.And

		case condition.IsElse(token):
			currentList = &newCondition.Else

		case condition.IsComparator(token):
			if newCondition.Target == "" {
				return errorCondition
			}
			newCondition.Comparator = condition.Comparator(token)

		default:
			push(token)
		}
	}

	if newCondition.Comparator == "" {
		return errorCondition
	}

	return *newCondition
}
