package utils

import "strings"

func IsStringStart(input string) bool {
	trimmed := strings.TrimSpace(input)

	if (strings.HasPrefix(trimmed, "\"") && len(trimmed) > 1) ||
		(strings.HasPrefix(trimmed, "'") && len(trimmed) > 1) {
		return true
	}

	return false
}

func IsStringEnd(input string) bool {
	trimmed := strings.TrimSpace(input)
	if len(trimmed) < 1 {
		return false
	}

	last := trimmed[len(trimmed)-1]
	if last != '"' && last != '\'' {
		return false
	}

	count := 0
	for i := len(trimmed) - 2; i >= 0 && trimmed[i] == '\\'; i-- {
		count++
	}

	return count%2 == 0
}
