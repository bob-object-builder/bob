package utils

import (
	"strconv"
	"strings"
)

func IsValue(input string) bool {
	trimmed := strings.TrimSpace(input)

	if len(trimmed) == 0 {
		return false
	}

	if (strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") && len(trimmed) > 1) ||
		(strings.HasPrefix(trimmed, "'") && strings.HasSuffix(trimmed, "'") && len(trimmed) > 1) {
		return true
	}

	if input == "?" {
		return true
	}

	if _, err := strconv.Atoi(trimmed); err == nil {
		return true
	}

	if _, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return true
	}

	lower := strings.ToLower(trimmed)
	if lower == "true" || lower == "false" {
		return true
	}

	if lower == "null" {
		return true
	}

	return false
}
