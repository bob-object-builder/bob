package lexer

import (
	"errors"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/utils"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"salvadorsru/bob/internal/lib/value/pill"
	"salvadorsru/bob/internal/models/get"
	"salvadorsru/bob/internal/models/insert"
	"salvadorsru/bob/internal/models/join"
	"salvadorsru/bob/internal/models/remove"
	"salvadorsru/bob/internal/models/table"
	"strings"
)

const CommentPrefix = "#"
const AliasSuffix = ":"

type Lexer struct {
	stack   Stack[any]
	tables  object.Object[table.Table]
	actions array.Array[any]
	pill    pill.Pill
	locked  bool
	jump    bool
	tokens  []string
	token   string
}

func New() *Lexer {
	return &Lexer{
		stack:   *NewStack[any](),
		tables:  *object.New[table.Table](),
		actions: *array.New[any](),
		pill:    pill.New(""),
		locked:  false,
		jump:    false,
		tokens:  []string{},
		token:   "",
	}
}

func (l *Lexer) Lock() {
	l.locked = true
}

func (l *Lexer) Unlock() {
	l.locked = false
}

func (l *Lexer) NextLine() {
	l.jump = true
}

func (l *Lexer) IsOpenKey() bool {
	return IsOpenKey(l.token)
}

func (l *Lexer) IsCloseKey() bool {
	return IsCloseKey(l.token)
}

func (l *Lexer) IsVoidContext() bool {
	return l.token == string(OpenKey)+string(CloseKey)
}

func (l *Lexer) GetTables() object.Object[table.Table] {
	return l.tables
}

func (l *Lexer) GetActions() array.Array[any] {
	return l.actions
}

func (l *Lexer) Parse(query string) (error, *object.Object[table.Table], *array.Array[any]) {

	capturing := false
	parametrising := false
	lines := utils.SplitByLine(query)

lineLoop:
	for line := range lines {

		tokens := strings.Fields(line.Content)
		l.tokens = tokens

		for _, token := range tokens {
			l.token = token

			if token == CommentPrefix {
				continue lineLoop
			}

			if l.IsOpenKey() {
				parametrising = false
				capturing = true
			}

			if l.IsCloseKey() {
				capturing = false
			}

			if l.jump {
				l.jump = false
				if !l.IsCloseKey() {
					continue lineLoop
				}
			}

			if checker.IsAlias(token) {
				aliasName := strings.TrimSuffix(token, AliasSuffix)
				l.pill.Set(aliasName)
				l.tokens = tokens[1:]
				continue
			}

			if parametrising {
				switch l.token {
				case table.Key, get.Key, join.LeftJoinKey, insert.Key, remove.Key:
					return errors.New("malformed query"), nil, nil
				}
			}

			switch token {
			case table.Key:
				l.stack.Push(table.New())
				parametrising = true
				continue
			case get.Key:
				l.stack.Push(get.New())
				parametrising = true
				continue
			case join.LeftJoinKey:
				l.stack.Push(join.NewLeftJoin())
				parametrising = true
				continue
			case insert.Key:
				l.stack.Push(insert.New())
				parametrising = true
				continue
			case remove.Key:
				l.stack.Push(remove.New())
				parametrising = true
				continue
			}

			item := l.stack.GetLast()

			if item == nil {
				continue
			}

			switch v := (*item).(type) {
			case *table.Table:
				if err := l.ParseTable(v); err != nil {
					return err, nil, nil
				}
			case *get.Get:
				l.ParseGet(v)
			case *join.Join:
				l.ParseLeftJoin(v)
			case *insert.Insert:
				if err := l.ParseInsert(v); err != nil {
					return err, nil, nil
				}
			case *remove.Remove:
				l.ParseRemove(v)
			}
		}
	}

	if parametrising || capturing {
		return errors.New("malformed query"), nil, nil
	}

	tables := l.GetTables()
	actions := l.GetActions()
	return nil, &tables, &actions
}
