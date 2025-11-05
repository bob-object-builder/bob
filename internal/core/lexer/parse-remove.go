package lexer

import (
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/remove"
)

func (l *Lexer) ParseRemove(r *remove.Remove) {
	if l.IsOpenKey() && !l.IsVoidContext() {
		return
	}

	if l.IsCloseKey() {
		l.actions.Push(*r)
		l.stack.Clean()
		return
	}

	if r.IsTargetEmpty() {
		r.SetTarget(l.token)
		return
	}

	if condition.IsCondition(l.token) {
		cond := l.ParseCondition(r.Target)
		r.Conditions.Push(cond)
		l.NextLine()
		return
	}
}
