package checker

import "unicode"

func StartWithUpperCase(s string) bool {
	return len(s) > 0 && unicode.IsUpper(rune(s[0]))
}
