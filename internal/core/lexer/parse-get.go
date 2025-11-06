package lexer

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/function"
	"salvadorsru/bob/internal/models/get"
)

func (l *Lexer) ParseGet(g *get.Get) {
	isVoidContext := l.IsVoidContext()

	if l.IsOpenKey() && !isVoidContext {
		g.Alias = l.pill.UseOr(g.Target)
		return
	}

	if l.IsCloseKey() || isVoidContext {
		if l.stack.Length() > 1 {
			l.stack.Pop()
			previous := l.stack.GetLast()

			if previousGet, ok := (*previous).(*get.Get); ok {
				previousGet.Subqueries.Push(*g)
			}

			return
		}

		l.actions.Push(*g)
		l.stack.Clean()
		return
	}

	if !g.HasTarget() {
		g.SetTarget(l.token)
		return
	}

	if condition.IsCondition(l.token) {
		cond := l.ParseCondition(g.Target)
		if g.Groups.Length() == 0 {
			g.Conditions.Push(cond)
		} else {
			g.Having.Push(cond)
		}
		l.NextLine()
		return
	}

	if get.IsGroup(l.token) {
		if len(l.tokens) > 1 {
			l.tokens = l.tokens[1:]
			g.Groups.Push(l.ParseReferences(g.Target))
		}
		l.NextLine()
		return
	}

	if function.IsFunction(l.token) {
		fn := l.ParseReferences(g.Target)
		g.Selected.Add(l.pill.UseOr(fn), fn)
		if len(l.tokens) > 1 {
			l.NextLine()
		}
		return
	}

	selected := formatter.ToSnakeCase(l.token)
	g.Selected.Add(l.pill.UseOr(selected), selected)
}
