package lexer

import (
	"salvadorsru/bob/internal/models/raw"
	"strings"
)

func (l *Lexer) ParseRaw(r *raw.Raw) error {
	if l.IsOpenKey() {
		return nil
	}

	if l.IsCloseKey() {
		l.actions.Push(*r)
		l.stack.Clean()
		return nil
	}

	if l.capturing {
		r.Lines.Push(strings.Join(l.tokens, " "))
	}
	l.NextLine()
	return nil
}
