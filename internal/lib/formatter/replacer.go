package formatter

import (
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/models/literal"
	"strings"
	"unicode"
)

// PrefixWith prepends `prefix` to valid identifiers in `target`.
//   - Reserved keywords (case-insensitive) are not modified (but see keywordFilter).
//   - Identifiers that are part of a dotted expression (e.g., user.orders) or arrow expression (e.g., user->orders)
//     are not modified, even if spaces exist around the dot or arrow.
//   - Words starting with '@' are considered literal and not modified.
//   - Content inside string literals ("string", 'char', `raw string`) is not modified,
//     except double quotes at the start are converted to single quotes.
//   - keywordFilter can return a replacement for a token, or empty string to ignore.
func PrefixWith(prefix string, target string, reservedKeywords []string, keywordFilter func(token string) string) string {
	isPrefixed := false
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

		// Handle identifiers and words starting with '@'
		if c == '@' || (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') || c == '_' {
			start := i
			if c == '@' {
				i++ // Skip '@' for word parsing
			}
			for i < n {
				ch := target[i]
				if (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_' {
					i++
				} else {
					break
				}
			}

			word := target[start:i]

			// If word starts with '@', treat it as literal
			if literal.IsLiteral(string(word[0])) {
				b.WriteString(literal.GetLiteral(word))
				continue
			}

			// Check reserved keywords
			if _, ok := reserved[word]; ok {
				filteredWord := keywordFilter(word)
				if filteredWord != "" {
					word = filteredWord
				}
				b.WriteString(word)
				continue
			}

			// Check if part of dotted or arrow expression
			prev := start - 1
			for prev >= 0 && unicode.IsSpace(rune(target[prev])) {
				prev--
			}
			if prev >= 0 {
				// Check '.' before
				if target[prev] == '.' {
					b.WriteString(word)
					continue
				}
				// Check '->' before
				if prev > 0 && target[prev-1] == '-' && target[prev] == '>' {
					b.WriteString(word)
					continue
				}
			}

			next := i
			for next < n && unicode.IsSpace(rune(target[next])) {
				next++
			}
			if next < n {
				// Check '.' after
				if target[next] == '.' {
					b.WriteString(word)
					continue
				}
				// Check '->' after
				if next+1 < n && target[next] == '-' && target[next+1] == '>' {
					b.WriteString(word)
					continue
				}
			}

			isPrefixed = true
			// Prefix the identifier
			b.WriteString(prefix)
			b.WriteString(word)

			continue
		}

		b.WriteByte(c)
		i++
	}

	token := b.String()

	if !isPrefixed && checker.IsReference(token) {
		token = ToReferenceCase(token)
	}

	return token
}
