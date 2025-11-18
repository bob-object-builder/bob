package lexer

import (
	"salvadorsru/bob/internal/models/raw"
	"strings"
)

func (l *Lexer) ParseRaw(r *raw.Raw) {
	if l.IsOpenKey() {
		return
	}

	if l.IsCloseKey() {
		l.actions.Push(*r)
		l.stack.Clean()
		return
	}

	if l.capturing {
		r.Lines.Push(strings.Join(l.tokens, " "))
	}
	l.NextLine()
}
