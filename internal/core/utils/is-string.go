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

	if (strings.HasSuffix(trimmed, "\"") && len(trimmed) > 1) ||
		(strings.HasSuffix(trimmed, "'") && len(trimmed) > 1) {
		return true
	}

	return false
}
