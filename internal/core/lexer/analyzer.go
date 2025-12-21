package lexer

import (
	"salvadorsru/bob/internal/core/drivers/driver"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
	"salvadorsru/bob/internal/model/drop"
	"salvadorsru/bob/internal/model/get"
	"salvadorsru/bob/internal/model/insert"
	"salvadorsru/bob/internal/model/join"
	"salvadorsru/bob/internal/model/raw"
	"salvadorsru/bob/internal/model/remove"
	"salvadorsru/bob/internal/model/set"
	"salvadorsru/bob/internal/model/table"
	"strings"
)

func (l *Lexer) onValue(token value.Token) *failure.Failure {
	word := token.Value

	isToken := func(check string) bool {
		return word == check
	}

	isNewLine := func() bool {
		return isToken(kw.NewLine)
	}

	switch word {
	case kw.OpenContext:
		l.isParameter = false
		return nil

	case kw.CloseContext, kw.VoidContext:
		if error := l.stack.Merge(); error != nil {
			return error
		}
		return nil

	case kw.Asc, kw.Desc:
		latest := *l.stack.GetLast()

		if g, ok := latest.(*get.Get); ok {
			l.stack.Push(get.NewOrder(g.Target, word))
		}

		return nil

	case kw.If, kw.Or:
		latest := *l.stack.GetLast()

		switch v := latest.(type) {
		case *get.Get:
			l.stack.Push(condition.NewCondition(v.Target, word))
		case *remove.Delete:
			l.stack.Push(condition.NewCondition(v.Target, word))
		case *join.Join:
			l.stack.Push(condition.NewCondition(v.ReferencedTable, word))
		case *set.Set:
			l.stack.Push(condition.NewCondition(v.Target, word))
		}

		return nil

	case kw.Group:
		latest := *l.stack.GetLast()

		if g, ok := latest.(*get.Get); ok {
			l.stack.Push(condition.NewGroup(g.Target))
		}

		return nil

	case kw.Table:
		l.stack.Push(table.NewTable())
		l.isParameter = true
		return nil

	case kw.Get:
		l.stack.Push(get.NewGet())
		l.isParameter = true
		return nil

	case kw.Join:
		l.stack.Push(join.NewJoin())
		l.isParameter = true
		return nil

	case kw.Delete:
		l.stack.Push(remove.NewDelete())
		l.isParameter = true
		return nil

	case kw.New:
		l.stack.Push(insert.NewInsert())
		l.isParameter = true
		return nil

	case kw.Set:
		l.stack.Push(set.NewSet())
		l.isParameter = true
		return nil

	case kw.Raw:
		l.stack.Push(raw.NewRaw())
		l.isParameter = true
		return nil

	case kw.Drop:
		l.stack.Push(drop.NewDrop())
		l.isParameter = true
		return nil
	}

	latest := l.stack.GetLast()

	if latest == nil {
		return nil
	}

	switch v := (*latest).(type) {
	case *table.Table:
		v.Parse(word)

		if isNewLine() {
			return nil
		}

		if !l.isParameter {
			if checker.IsUppercase(word) {
				l.stack.Push(table.NewReference(word))
				return nil
			} else {
				column := table.NewColumn(word)
				l.stack.Push(column)
				return nil
			}
		}

	case *table.Column:
		if err := v.Parse(word); err != nil {
			return err
		}

		if isNewLine() {
			if error := l.stack.Merge(); error != nil {
				return error
			}
		}

	case *table.Reference:
		v.Parse(word)

		if isNewLine() {
			if error := l.stack.Merge(); error != nil {
				return error
			}
		}

	case *get.Order:
		if isNewLine() {
			if error := l.stack.Merge(); error != nil {
				return error
			}
		}

		if err := v.Parse(word); err != nil {
			return err
		}

	case *condition.Condition:
		if isNewLine() {
			if error := l.stack.Merge(); error != nil {
				return error
			}
		}

		v.Parse(token)

	case *condition.Group:
		if isNewLine() {
			if error := l.stack.Merge(); error != nil {
				return error
			}
		}

		v.Parse(word)

	case *get.Get:
		if l.isParameter {
			v.Parse(word)
			return nil
		}

		if isNewLine() {
			return nil
		}

		if token.Type == value.Alias || token.Type == value.Key {
			l.stack.Push(get.NewGetChild(v.Target, token.Value))
			return nil
		} else {
			return failure.CannotUseExpressionOrStringAsSelection
		}

	case *get.Child:
		if isNewLine() {
			if error := l.stack.Merge(); error != nil {
				return error
			}
		}

		v.Parse(token)

	case *join.Join:
		if l.isParameter {
			v.Parse(word, l.isParameter)
			return nil
		}

		if isNewLine() {
			return nil
		}

		if token.Type == value.Alias || token.Type == value.Key {
			l.stack.Push(join.NewChild(v.ReferencedTable, token))
			return nil
		} else {
			return failure.CannotUseExpressionOrStringAsSelection
		}

	case *join.Selection:
		if isNewLine() {
			if error := l.stack.Merge(); error != nil {
				return error
			}
		}

		if error := v.Parse(token); error != nil {
			return error
		}

	case *remove.Delete:
		v.Parse(token)

	case *insert.Insert:
		v.Parse(token, l.isParameter)

	case *raw.Raw:
		if !l.isParameter {
			v.Parse(token)
		}

	case *drop.Drop:
		v.Parse(token)

	case *set.Set:
		if error := v.Parse(token); error != nil {
			return error
		}

	}

	return nil
}

func (l *Lexer) onExpression(token string) string {
	lt := *l.stack.GetLast()

	prefix := ""

	switch v := lt.(type) {
	case *get.Get:
		prefix = v.Target
	case *set.Set:
		prefix = v.Target
	case *get.Order:
		prefix = v.From
	case *get.Child:
		prefix = v.From
	case *join.Selection:
		prefix = v.From
	case *condition.Condition:
		prefix = v.From
	}

	if kw.IsFunction(token) {
		return strings.ToUpper(token)
	}

	if checker.IsNumber(token) {
		return token
	}

	if driver.IsLiteral(token) {
		return driver.GetLiteral(token)
	}

	if prefix == "" {
		return token
	}

	return prefix + "." + token
}
