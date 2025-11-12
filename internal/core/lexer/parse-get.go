package lexer

import (
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/function"
	"salvadorsru/bob/internal/models/get"
	"strings"
)

func (l *Lexer) ParseGet(g *get.Get) *failure.Failure {
	isVoidContext := l.IsVoidContext()

	if l.IsOpenKey() && !isVoidContext {
		g.Alias = l.pill.UseOr(g.Target)
		return nil
	}

	if l.IsCloseKey() || isVoidContext {
		if l.stack.Length() > 1 {
			l.stack.Pop()
			previous := l.stack.GetLast()

			if previousGet, ok := (*previous).(*get.Get); ok {
				previousGet.Subqueries.Push(*g)
			}

			return nil
		}

		l.actions.Push(*g)
		l.stack.Clean()
		return nil
	}

	if !g.HasTarget() {
		g.SetTarget(l.token)
		return nil
	}

	if get.IsLimit(l.token) {
		tokens := l.tokens

		if len(tokens) < 2 {
			return failure.LimitValueRequired
		}

		limitValue := tokens[1]
		if !checker.IsInt(limitValue) {
			return failure.LimitValueMustBeInteger
		}
		g.Limit = limitValue

		if len(tokens) >= 4 && get.IsOffset(tokens[2]) {
			offsetValue := tokens[3]
			if !checker.IsInt(offsetValue) {
				return failure.OffsetValueMustBeInteger
			}
			g.Offset = offsetValue
		}

		l.NextLine()
		return nil
	}

	if condition.IsCondition(l.token) {
		cond := l.ParseCondition(g.Target)
		if g.Groups.Length() == 0 {
			g.Conditions.Push(cond)
		} else {
			g.Having.Push(cond)
		}
		l.NextLine()
		return nil
	}

	if get.IsGroup(l.token) {
		if len(l.tokens) > 1 {
			l.tokens = l.tokens[1:]
			g.Groups.Push(l.ParseReferences(g.Target))
		}
		l.NextLine()
		return nil
	}

	if function.IsFunction(l.token) {
		fn := l.ParseReferences(g.Target)
		g.Selected.Add(l.pill.UseOr(fn), fn)
		if len(l.tokens) > 1 {
			l.NextLine()
		}
		return nil
	}

	selected := l.token

	if strings.Contains(l.token, ".") {
		selected = formatter.ToSnakeCase(selected)
	}

	if !checker.IsWord(selected) {
		return failure.InvalidSelectedColumn(selected)
	}

	g.Selected.Add(l.pill.UseOr(selected), selected)
	return nil
}
