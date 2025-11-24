package lexer

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/set"
)

func (l *Lexer) ParseSet(s *set.Set) *failure.Failure {
	if l.IsOpenKey() && !l.IsVoidContext() {
		return nil
	}

	if l.IsCloseKey() {
		l.actions.Push(*s)
		l.stack.Clean()
		return nil
	}

	if s.IsTargetEmpty() {
		s.SetTarget(l.token)
		return nil
	}

	if condition.IsCondition(l.token) {
		cond := l.ParseCondition(s.Target)
		s.Conditions.Push(cond)
		l.NextLine()
		return nil
	}

	if !checker.IsWord(l.token) {
		return failure.InvalidSetter
	}

	l.tokens = l.tokens[1:]
	s.Values.Add(l.token, l.ParseReferences(s.Target))
	l.NextLine()
	return nil
}
