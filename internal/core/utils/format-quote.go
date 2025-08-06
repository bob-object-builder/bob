package utils

import "strings"

func FormatQuote(s string) string {
	if len(s) >= 2 && strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`) {
		return "'" + s[1:len(s)-1] + "'"
	}
	return s
}
