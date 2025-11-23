package lexer

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/insert"
	"strings"
)

func (l *Lexer) ParseInsert(i *insert.Insert) *failure.Failure {
	if l.IsOpenKey() {
		if i.Columns.Length() > 0 {
			i.IsBulk = true
		}
		i.Capturing = true
		return nil
	}

	if l.IsCloseKey() {
		l.actions.Push(*i)
		l.stack.Clean()
		return nil
	}

	if i.IsTargetEmpty() {
		i.SetTarget(formatter.ToSnakeCase(l.token))
		return nil
	}

	if !i.Capturing {
		if strings.Contains(l.token, ".") {
			i.AddColumn(formatter.ToSnakeCase(l.token))
		} else {
			i.AddColumn(l.token)
		}
		return nil
	}

	if !i.IsBulk {
		i.AddColumn(l.token)
		l.tokens = l.tokens[1:]
	}

	buffering := false
	buffer := array.New[string]()
	values := array.New[string]()

	for _, token := range l.tokens {

		if buffering || checker.IsStringStart(token) || checker.IsExpressionStart(token) {
			buffering = true
			buffer.Push(token)

			if checker.IsExpressionEnd(token) || checker.IsStringEnd(token) {
				buffering = false
				joined := strings.Join(*buffer, " ")
				buffer.Clean()
				values.Push(formatter.NormalizeString(joined))
			}

			continue
		}

		values.Push(token)
	}

	if i.IsBulk {
		if i.Columns.Length() != values.Length() {
			return failure.ColumnNotReceivingValue(*i.Columns.GetLast(), strings.Join(l.tokens, ", "))
		}

		i.Values.Push(*values)
	} else {
		row := i.Values.Get(0)
		if row == nil {
			empty := array.New[string]()
			i.Values.Push(*empty)
			row = i.Values.Get(0)
		}
		row.Push(strings.Join(*values, " "))
	}

	l.NextLine()

	return nil
}
