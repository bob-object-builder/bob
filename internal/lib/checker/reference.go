package checker

import "strings"

func IsReference(s string) bool {
	return (strings.Contains(s, ".") || strings.Contains(s, "->")) && !IsInt(s) && !IsStringStart(s)
}
