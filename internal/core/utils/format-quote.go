package utils

import "strings"

func FormatQuote(s string) string {
	if len(s) >= 2 && strings.HasPrefix(s, `"`) {
		s = "'" + s[1:]
	}

	if len(s) >= 2 && strings.HasSuffix(s, `"`) {
		s = s[0:len(s)-1] + "'"
	}

	return s
}
