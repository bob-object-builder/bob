package formatter

import (
	"strings"
	"unicode"
)

func ToSnakeCase(str string) string {
	var result strings.Builder
	runes := []rune(str)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 {
				prev := runes[i-1]
				if unicode.IsLower(prev) || unicode.IsDigit(prev) ||
					(i+1 < len(runes) && unicode.IsLower(runes[i+1]) && unicode.IsUpper(prev)) {
					result.WriteRune('_')
				}
			}
			result.WriteRune(unicode.ToLower(r))
		} else if r == ' ' || r == '.' {
			result.WriteRune('_')
		} else {
			result.WriteRune(r)
		}
	}

	return result.String()
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

func String(s string) string {
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
