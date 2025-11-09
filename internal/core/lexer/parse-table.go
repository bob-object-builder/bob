package lexer

import (
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/failure"
	"salvadorsru/bob/internal/models/table"
)

func (l *Lexer) ParseTable(t *table.Table) *failure.Failure {
	if l.IsOpenKey() {
		return nil
	}

	if l.IsCloseKey() {
		l.tables.Add(t.GetName(), *t)
		l.stack.Clean()
		return nil
	}

	if t.IsNameEmpty() {
		t.SetName(l.token)
		return nil
	}

	if checker.StartWithUpperCase(l.token) {
		table := l.token

		column := "id"
		if len(l.tokens) > 1 {
			column = l.tokens[1]
		}

		properties := []string{}
		if len(l.tokens) > 2 {
			properties = l.tokens[2:]
		}
		err := t.AddReference(table, column, properties)
		if err != nil {
			return err
		}
	} else {
		if len(l.tokens) < 2 {
			return failure.UndefinedToken(l.token)
		}

		name, kind, properties := l.token, l.tokens[1], l.tokens[2:]
		err := t.AddColumn(name, kind, properties)
		if err != nil {
			return err
		}
	}

	l.NextLine()
	return nil
}
