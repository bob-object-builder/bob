package formatter

import "strings"

func ToReferenceCase(str string) string {
	parts := strings.Split(str, "->")

	for i, part := range parts {
		if len(part) > 0 && part[0] >= 'A' && part[0] <= 'Z' {
			part = strings.ToLower(part[:1]) + part[1:]
		}
		part = strings.ReplaceAll(part, ".", "_")
		parts[i] = part
	}
	return strings.Join(parts, ".")
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
