package formatter

import (
	"strings"
	"unicode"
)

// PrefixWith prepends `prefix` to valid identifiers in `target`.
//   - Reserved keywords (case-insensitive) are not modified.
//   - Identifiers that are part of a dotted expression (e.g., user.orders)
//     are not modified, even if spaces exist around the dot.
//   - Content inside string literals ("string", 'char', `raw string`) is not modified,
//     except double quotes at the start are converted to single quotes.
func PrefixWith(prefix string, target string, reservedKeywords []string) string {
	reserved := make(map[string]struct{}, len(reservedKeywords))
	for _, kw := range reservedKeywords {
		reserved[kw] = struct{}{}
	}

	var b strings.Builder
	n := len(target)
	i := 0

	for i < n {
		c := target[i]

		// Handle string/char/raw literals
		if c == '"' || c == '\'' || c == '`' {
			startQuote := c
			if c == '"' {
				// Replace starting double quote with single quote
				b.WriteByte('\'')
			} else {
				b.WriteByte(c)
			}
			i++
			for i < n {
				ch := target[i]
				if ch == '\\' && startQuote != '`' {
					// Escape sequence: copy both characters
					b.WriteByte(ch)
					i++
					if i < n {
						b.WriteByte(target[i])
						i++
					}
					continue
				}
				if ch == startQuote {
					// End of literal
					if startQuote == '"' {
						b.WriteByte('\'') // Replace ending double quote
					} else {
						b.WriteByte(ch)
					}
					i++
					break
				}
				b.WriteByte(ch)
				i++
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

			if _, ok := reserved[word]; ok {
				b.WriteString(word)
				continue
			}

			prev := start - 1
			for prev >= 0 && unicode.IsSpace(rune(target[prev])) {
				prev--
			}
			if prev >= 0 && target[prev] == '.' {
				b.WriteString(word)
				continue
			}

			next := i
			for next < n && unicode.IsSpace(rune(target[next])) {
				next++
			}
			if next < n && target[next] == '.' {
				b.WriteString(word)
				continue
			}

			b.WriteString(prefix)
			b.WriteString(word)
			continue
		}

		b.WriteByte(c)
		i++
	}

	return b.String()
}
