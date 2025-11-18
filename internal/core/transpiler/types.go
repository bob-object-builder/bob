package transpiler

import (
	"salvadorsru/bob/internal/lib/value/array"
	"strings"
)

type TranspiledTable struct {
	table array.Array[string]
}

func (t *TranspiledTable) ToString() string {
	return strings.Join(t.table, "\n\n")
}

func (t *TranspiledTable) Push(value ...string) {
	t.table.Push(value...)
}

func (t *TranspiledTable) Get() []string {
	return t.table
}

type TranspiledActions struct {
	actions array.Array[string]
}

func (t *TranspiledActions) ToString() string {
	return strings.Join(t.actions, "\n\n")
}

func (t *TranspiledActions) Push(value ...string) {
	t.actions.Push(value...)
}

func (t *TranspiledActions) Get() []string {
	return t.actions
}
