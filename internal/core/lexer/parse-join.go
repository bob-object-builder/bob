package lexer

import (
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/get"
	"salvadorsru/bob/internal/models/join"
)

func (l *Lexer) ParseLeftJoin(j *join.Join) {
	if l.IsOpenKey() {
		j.Capturing = false

		if j.IsOnEmpty() {
			j.SetOn("id")
		}

		return
	}

	if l.IsCloseKey() {
		if l.stack.Length() > 1 {
			l.stack.Pop()
			previous := l.stack.GetLast()

			switch previous := (*previous).(type) {
			case *get.Get:
				previous.Joins.Push(*j)
			case *join.Join:
				previous.Subjoins.Push(*j)
			}

			return
		}

		l.stack.Clean()
		return
	}

	if j.IsTargetEmpty() {
		j.SetTarget(l.token)
		j.Capturing = false
		return

	}

	if j.IsOnEmpty() {
		j.SetOn(l.token)
		return
	}

	if condition.IsCondition(l.token) {
		cond := l.ParseCondition(j.Target)
		j.Conditions.Push(cond)
		l.NextLine()
		return
	}

	if get.IsGroup(l.token) {
		if len(l.tokens) > 1 {
			l.tokens = l.tokens[1:]
			j.Groups.Push(l.ParseReferences(j.Target))
		}
		l.NextLine()
		return
	}

	j.Selected.Add(l.pill.UseOr(l.token), l.ParseReferences(j.Target))
	if len(l.tokens) > 1 {
		l.NextLine()
	}
}
