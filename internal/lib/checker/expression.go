package checker

import "strings"

func IsExpressionStart(input string) bool {
	trimmed := strings.TrimSpace(input)
	return len(trimmed) > 0 && (strings.HasPrefix(trimmed, "("))
}

func IsExpressionEnd(input string) bool {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) < 1 {
		return false
	}
	last := trimmed[len(trimmed)-1]
	return last == ')'
}
