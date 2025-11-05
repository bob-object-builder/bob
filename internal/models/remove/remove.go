package remove

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/condition"
)

const Key = "delete"

type Remove struct {
	Target     string
	Conditions array.Array[condition.Condition]
}

func New() *Remove {
	return &Remove{}
}

func (get *Remove) IsTargetEmpty() bool {
	return get.Target == ""
}

func (get *Remove) SetTarget(target string) {
	get.Target = formatter.ToSnakeCase(target)
}
