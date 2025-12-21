package lexer

import (
	"strings"

	"salvadorsru/bob/internal/core/drivers/driver"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/stack"
	"salvadorsru/bob/internal/lib/value"
)

type Lexer struct {
	Driver      driver.Driver
	isParameter bool
	stack       stack.Stack[any]
}

var latestKey string

func (l *Lexer) Parse(query string) (*failure.Failure, *stack.Stack[any]) {

	var word string

	inString := false
	var delimiter rune
	escape := false

	inExpr := false
	depth := 0
	expr_buffer := ""

	// =========================
	//   TOKEN EMITTERS
	// =========================
	emitKey := func(v string) *failure.Failure {
		if latestKey == kw.Raw {
			return l.onValue(value.Token{
				Type:  value.Key,
				Value: v,
			})
		}

		if kw.IsKeyword(v) {
			latestKey = v
			return l.onValue(value.Token{
				Type:  value.Key,
				Value: v,
			})
		}

		return l.onValue(value.Token{
			Type:  value.Key,
			Value: v,
		})
	}

	emitText := func(v string) *failure.Failure {
		return l.onValue(value.Token{
			Type:  value.Text,
			Value: v,
		})
	}

	emitAlias := func(v string) *failure.Failure {
		return l.onValue(value.Token{
			Type:  value.Alias,
			Value: v,
		})
	}

	emitExpression := func(v string) *failure.Failure {
		return l.onValue(value.Token{
			Type:  value.Expression,
			Value: v,
		})
	}

	// =========================
	//   VALUE EMITTER
	// =========================
	emitValue := func() *failure.Failure {
		if inExpr {
			if len(word) > 0 {
				processed := l.onExpression(word)
				if processed == "" {
					return failure.ExpressionProcessingEmptyResult
				}
				expr_buffer += processed
			}
			word = ""
			return nil
		}

		word = strings.TrimSpace(word)
		if len(word) == 0 {
			return nil
		}

		if checker.IsAlias(word) {
			if err := emitAlias(strings.TrimSuffix(word, ":")); err != nil {
				return err
			}
		} else {
			if err := emitKey(word); err != nil {
				return err
			}
		}

		word = ""
		return nil
	}

	// =========================
	//   MAIN LOOP (UTF-8 SAFE)
	// =========================
	for _, ch := range query {

		// ========================
		//   STRING MODE
		// ========================
		if inString {
			if escape {
				word += string(ch)
				escape = false
				continue
			}

			if ch == '\\' {
				word += string(ch)
				escape = true
				continue
			}

			if ch == delimiter {
				word += string(ch)

				// normalize " → '
				if delimiter == '"' && len(word) >= 2 {
					word = "'" + word[1:len(word)-1] + "'"
				}

				// escape internal apostrophes
				if len(word) >= 2 && word[0] == '\'' && word[len(word)-1] == '\'' {
					inner := word[1 : len(word)-1]
					inner = strings.ReplaceAll(inner, "'", "''")
					word = "'" + inner + "'"
				}

				if inExpr {
					expr_buffer += word
				} else {
					if err := emitText(word); err != nil {
						return err, nil
					}
				}

				word = ""
				inString = false
				continue
			}

			word += string(ch)
			continue
		}

		// =========================
		//   ENTER STRING
		// =========================
		if ch == '"' || ch == '\'' {
			if err := emitValue(); err != nil {
				return err, nil
			}
			inString = true
			delimiter = ch
			word = string(ch)
			continue
		}

		// =======================
		//   EXPRESSION MODE
		// =======================
		if inExpr {
			switch ch {
			case '(':
				if err := emitValue(); err != nil {
					return err, nil
				}
				depth++
				expr_buffer += "("

			case ')':
				if depth == 0 {
					return failure.UnexpectedClosingParenthesis, nil
				}

				if err := emitValue(); err != nil {
					return err, nil
				}

				expr_buffer += ")"
				depth--

				if depth == 0 {
					inExpr = false
					if err := emitExpression(expr_buffer); err != nil {
						return err, nil
					}
					expr_buffer = ""
				}

			case ',':
				if err := emitValue(); err != nil {
					return err, nil
				}
				expr_buffer += ","

			case ' ', '\n', '\t':
				if len(word) > 0 {
					expr_buffer += word
					word = ""
				}
				expr_buffer += string(ch)

			default:
				word += string(ch)
			}
			continue
		}

		// =======================
		//  OUTSIDE EXPRESSION
		// =======================
		switch ch {
		case '(':
			if len(word) > 0 {
				if !kw.IsFunction(word) {
					return failure.UndefinedFunction(word), nil
				}

				processed := l.onExpression(word)
				if processed == "" {
					return failure.InvalidExpressionBeforeParenthesis, nil
				}
				expr_buffer += processed
				word = ""
			}
			inExpr = true
			depth = 1
			expr_buffer += "("

		case ')':
			if err := emitValue(); err != nil {
				return err, nil
			}
			if err := emitKey(")"); err != nil {
				return err, nil
			}

		case ',':
			if err := emitValue(); err != nil {
				return err, nil
			}
			if err := emitKey(","); err != nil {
				return err, nil
			}

		case ' ', '\t':
			if err := emitValue(); err != nil {
				return err, nil
			}

		case '\n':
			if err := emitValue(); err != nil {
				return err, nil
			}
			if err := emitKey("\n"); err != nil {
				return err, nil
			}

		case '{', '}':
			if err := emitValue(); err != nil {
				return err, nil
			}
			if err := emitKey(string(ch)); err != nil {
				return err, nil
			}

		default:
			word += string(ch)
		}
	}

	if err := emitValue(); err != nil {
		return err, nil
	}

	if expr_buffer != "" {
		if inExpr {
			return failure.UnclosedExpression, nil
		}
		if err := emitExpression(expr_buffer); err != nil {
			return err, nil
		}
	}

	if inString {
		return failure.UnclosedStringLiteral, nil
	}

	return nil, &l.stack
}
