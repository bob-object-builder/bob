package formatter

import (
	"strings"
	"unicode"
)

// PrefixWith prepends `prefix` to valid identifiers in `target`.
//   - Reserved keywords (case-insensitive) are not modified.
//   - Identifiers that are part of a dotted expression (e.g., user.orders)
//     are not modified, even if spaces exist around the dot.
//   - Content inside string literals ("string", 'char', `raw string`) is not modified.
func PrefixWith(prefix string, target string, reservedKeywords []string) string {
	// Prepare a set of reserved keywords in lowercase for case-insensitive comparison
	reserved := make(map[string]struct{}, len(reservedKeywords))
	for _, kw := range reservedKeywords {
		reserved[strings.ToLower(kw)] = struct{}{}
	}

	var b strings.Builder
	n := len(target)
	i := 0

	for i < n {
		c := target[i]

		// Handle string/char/raw literals
		if c == '"' || c == '\'' || c == '`' {
			startQuote := c
			b.WriteByte(c)
			i++
			for i < n {
				ch := target[i]
				b.WriteByte(ch)
				i++
				if ch == '\\' && startQuote != '`' {
					// Escape sequence: copy the next character if it exists
					if i < n {
						b.WriteByte(target[i])
						i++
					}
					continue
				}
				if ch == startQuote {
					break
				}
			}
			continue
		}

		// Start of a possible identifier (ASCII letters + underscore)
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_' {
			start := i
			for i < n {
				ch := target[i]
				if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' {
					i++
				} else {
					break
				}
			}

			word := target[start:i]
			lower := strings.ToLower(word)

			// Reserved keyword → leave unchanged
			if _, ok := reserved[lower]; ok {
				b.WriteString(word)
				continue
			}

			// Check previous non-space character
			prev := start - 1
			for prev >= 0 && unicode.IsSpace(rune(target[prev])) {
				prev--
			}
			if prev >= 0 && target[prev] == '.' {
				b.WriteString(word)
				continue
			}

			// Check next non-space character
			next := i
			for next < n && unicode.IsSpace(rune(target[next])) {
				next++
			}
			if next < n && target[next] == '.' {
				b.WriteString(word)
				continue
			}

			// Valid identifier → prepend prefix
			b.WriteString(prefix)
			b.WriteString(word)
			continue
		}

		// Any other character: copy as-is
		b.WriteByte(c)
		i++
	}

	return b.String()
}
