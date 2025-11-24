package lexer

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/get"
	"salvadorsru/bob/internal/models/join"
	"salvadorsru/bob/internal/models/order"
)

func (l *Lexer) ParseLeftJoin(j *join.Join) *failure.Failure {
	if l.IsOpenKey() {
		j.Capturing = false

		if j.IsOnEmpty() {
			j.SetOn("id")
		}

		return nil
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

			return nil
		}

		l.stack.Clean()
		return nil
	}

	if j.IsTargetEmpty() {
		j.SetTarget(l.token)
		j.Capturing = false
		return nil

	}

	if j.IsOnEmpty() {
		j.SetOn(l.token)
		return nil
	}

	if condition.IsCondition(l.token) {
		cond := l.ParseCondition(j.Target)
		j.Conditions.Push(cond)
		l.NextLine()
		return nil
	}

	if get.IsGroup(l.token) {
		if len(l.tokens) > 1 {
			l.tokens = l.tokens[1:]
			j.Groups.Push(l.ParseReferences(j.Target))
		}
		l.NextLine()
		return nil
	}

	if order.IsOrder(l.token) {
		orderError, order := l.ParseOrder(j.Target)
		if orderError != nil {
			return orderError
		}

		j.Orders.Push(*order)
		return nil
	}

	j.Selected.Add(l.pill.UseOr(l.token), l.ParseReferences(j.Target))
	if len(l.tokens) > 1 {
		l.NextLine()
	}

	return nil
}
