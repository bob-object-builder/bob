package lexer

import (
	"salvadorsru/bob/internal/models/drop"
)

func (l *Lexer) ParseDrop(d *drop.Drop) error {
	isVoidContext := l.IsVoidContext()

	if l.IsOpenKey() && !isVoidContext {
		return nil
	}

	if l.IsCloseKey() || l.IsVoidContext() {
		l.actions.Push(*d)
		l.stack.Clean()
		return nil
	}

	if !d.HasTarget() {
		d.SetTarget(l.token)
		return nil
	}

	return nil
}
