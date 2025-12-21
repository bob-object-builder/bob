package table

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/value"
)

type Table struct {
	Name       string
	Columns    value.Map[Column]
	References value.Map[Reference]
}

func NewTable() *Table {
	return &Table{}
}

func (t *Table) Merge(i any) *failure.Failure {

	switch v := i.(type) {
	case *Column:
		t.Columns.Set(v.Name, *v)
	case *Reference:
		t.References.Set(v.Target, *v)
	}

	return nil
}

func (t *Table) Parse(token string) {
	if t.Name == "" {
		t.Name = token
		return
	}
}
