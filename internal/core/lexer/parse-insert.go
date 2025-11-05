package lexer

import (
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/insert"
	"strings"
)

func (l *Lexer) ParseInsert(i *insert.Insert) {
	if l.IsOpenKey() {
		if i.Columns.Length() > 0 {
			i.IsBulk = true
		}
		i.Capturing = true
		return
	}

	if l.IsCloseKey() {
		l.actions.Push(*i)
		l.stack.Clean()
		return
	}

	if i.IsTargetEmpty() {
		i.SetTarget(l.token)
		return
	}

	if !i.Capturing {
		i.AddColumn(l.token)
		return
	}

	if !i.IsBulk {
		i.AddColumn(l.token)
		l.tokens = l.tokens[1:]
	}

	parsing_string := false
	buffer := array.New[string]()
	values := array.New[string]()

	for _, token := range l.tokens {
		if parsing_string || checker.IsStringStart(token) {
			parsing_string = true
			buffer.Push(token)

			if checker.IsStringEnd(token) {
				parsing_string = false
				joined := strings.Join(*buffer, " ")
				buffer.Clean()
				values.Push(formatter.String(joined))
			}
			continue
		}

		values.Push(token)
	}

	if i.IsBulk {
		i.Values.Push(*values)
	} else {
		row := i.Values.Get(0)
		if row == nil {
			empty := array.New[string]()
			i.Values.Push(*empty)
			row = i.Values.Get(0)
		}
		row.Push(*values...)
	}

	l.NextLine()
}
