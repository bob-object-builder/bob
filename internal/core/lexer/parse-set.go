package lexer

import (
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/set"
)

func (l *Lexer) ParseSet(s *set.Set) {
	if l.IsOpenKey() && !l.IsVoidContext() {
		return
	}

	if l.IsCloseKey() {
		l.actions.Push(*s)
		l.stack.Clean()
		return
	}

	if s.IsTargetEmpty() {
		s.SetTarget(l.token)
		return
	}

	if condition.IsCondition(l.token) {
		cond := l.ParseCondition(s.Target)
		s.Conditions.Push(cond)
		l.NextLine()
		return
	}

	l.tokens = l.tokens[1:]
	s.Values.Add(l.token, l.ParseReferences(s.Target))
	l.NextLine()
}
