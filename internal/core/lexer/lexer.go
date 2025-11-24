package lexer

import (
	"salvadorsru/bob/internal/core/driver"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/utils"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"salvadorsru/bob/internal/lib/value/pill"
	"salvadorsru/bob/internal/models/drop"
	"salvadorsru/bob/internal/models/get"
	"salvadorsru/bob/internal/models/insert"
	"salvadorsru/bob/internal/models/join"
	"salvadorsru/bob/internal/models/raw"
	"salvadorsru/bob/internal/models/remove"
	"salvadorsru/bob/internal/models/set"
	"salvadorsru/bob/internal/models/table"
	"strings"
)

const CommentPrefix = "#"
const AliasSuffix = ":"

type Lexer struct {
	stack         Stack[any]
	tables        object.Object[table.Table]
	actions       array.Array[any]
	pill          pill.Pill
	token         string
	tokens        []string
	locked        bool
	jump          bool
	capturing     bool
	parametrising bool
	driver        driver.Driver
}

func New(drv driver.Driver) *Lexer {
	return &Lexer{
		stack:         *NewStack[any](),
		tables:        *object.New[table.Table](),
		actions:       *array.New[any](),
		pill:          pill.New(""),
		locked:        false,
		jump:          false,
		tokens:        []string{},
		token:         "",
		capturing:     false,
		parametrising: false,
		driver:        drv,
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
	hasOpenKey := strings.HasPrefix(l.token, string(OpenKey))
	hasCloseKey := strings.HasSuffix(l.token, string(CloseKey))
	return hasOpenKey && hasCloseKey
}

func (l *Lexer) GetTables() object.Object[table.Table] {
	return l.tables
}

func (l *Lexer) GetActions() array.Array[any] {
	return l.actions
}

func (l *Lexer) Parse(query string) (*failure.Failure, *object.Object[table.Table], *array.Array[any]) {
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
				l.parametrising = false
				l.capturing = true
			}

			if l.IsCloseKey() {
				l.capturing = false
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

			if l.parametrising {
				switch l.token {
				case table.Key, get.Key, join.LeftJoinKey, insert.Key, remove.Key:
					return failure.MalformedQuery(l.token), nil, nil
				}
			}

			switch token {
			case table.Key:
				l.stack.Push(table.New())
				l.parametrising = true
				continue
			case get.Key:
				l.stack.Push(get.New())
				l.parametrising = true
				continue
			case join.LeftJoinKey:
				l.stack.Push(join.NewLeftJoin())
				l.parametrising = true
				continue
			case insert.Key:
				l.stack.Push(insert.New())
				l.parametrising = true
				continue
			case remove.Key:
				l.stack.Push(remove.New())
				l.parametrising = true
				continue
			case raw.Key:
				l.stack.Push(raw.New())
				l.parametrising = true
				continue
			case drop.Key:
				l.stack.Push(drop.New())
				l.parametrising = true
				continue
			case set.Key:
				l.stack.Push(set.New())
				l.parametrising = true
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
				if err := l.ParseGet(v); err != nil {
					return err, nil, nil
				}
			case *join.Join:
				l.ParseLeftJoin(v)
			case *insert.Insert:
				if err := l.ParseInsert(v); err != nil {
					return err, nil, nil
				}
			case *remove.Remove:
				l.ParseRemove(v)
			case *raw.Raw:
				l.ParseRaw(v)
			case *drop.Drop:
				l.ParseDrop(v)
			case *set.Set:
				l.ParseSet(v)
			}
		}
	}

	if l.parametrising || l.capturing {
		return failure.MalformedQuery(l.token), nil, nil
	}

	tables := l.GetTables()
	actions := l.GetActions()
	return nil, &tables, &actions
}
