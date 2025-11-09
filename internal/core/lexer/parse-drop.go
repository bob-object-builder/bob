package lexer

import (
	"salvadorsru/bob/internal/models/drop"
)

func (l *Lexer) ParseDrop(d *drop.Drop) {
	isVoidContext := l.IsVoidContext()

	if l.IsOpenKey() && !isVoidContext {
		return
	}

	if l.IsCloseKey() || l.IsVoidContext() {
		l.actions.Push(*d)
		l.stack.Clean()
		return
	}

	if !d.HasTarget() {
		d.SetTarget(l.token)
		return
	}
}
