package insert

import (
	"salvadorsru/bob/internal/lib/value/array"
)

const Key = "new"

type Insert struct {
	Target    string
	Columns   array.Array[string]
	Values    array.Array[array.Array[string]]
	IsBulk    bool
	Capturing bool
}

func New() *Insert {
	return &Insert{Capturing: false}
}

func (i *Insert) IsTargetEmpty() bool {
	return i.Target == ""
}

func (i *Insert) SetTarget(target string) {
	i.Target = target
}

func (i *Insert) AddColumn(column string) {
	i.Columns.Push(column)
}
