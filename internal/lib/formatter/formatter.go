package formatter

import (
	"strings"
)

func ToReferenceCase(str string) string {
	if len(str) > 0 && str[0] >= 'A' && str[0] <= 'Z' {
		str = strings.ToLower(str[:1]) + str[1:]
	}

	return strings.ReplaceAll(str, ".", "_")
}

func Indent(str string, size ...int) string {
	indentSize := 1
	if len(size) > 0 {
		indentSize = size[0]
	}
	indent := strings.Repeat("  ", indentSize)
	return indent + str
}

func IndentLines(str string, size ...int) string {
	indentSize := 1
	if len(size) > 0 {
		indentSize = size[0]
	}
	indent := strings.Repeat("  ", indentSize)
	lines := strings.Split(str, "\n")
	for i, line := range lines {
		lines[i] = indent + line
	}
	return strings.Join(lines, "\n")
}

func NormalizeString(s string) string {
	if s == "" {
		return s
	}

	if s[0] == '"' {
		s = "'" + s[1:]
	}

	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1] + "'"
	}

	return s
}
