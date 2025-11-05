package lexer

import (
	"fmt"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/models/table"
)

func (l *Lexer) ParseTable(t *table.Table) error {
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
		name := l.token

		column := "id"
		if len(l.tokens) > 1 {
			column = l.tokens[1]
		}

		isOptional := false
		if len(l.tokens) > 2 {
			isOptional = l.tokens[2] == string(table.OptionalKey)
		}

		t.AddReference(name, column, isOptional)
	} else {
		if len(l.tokens) < 2 {
			return fmt.Errorf("undefined token '%s'", l.token)
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
