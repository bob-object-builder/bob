package utils

import "strings"

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
