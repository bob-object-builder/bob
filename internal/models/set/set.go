package set

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"salvadorsru/bob/internal/models/condition"
)

const Key = "set"

type Set struct {
	Target     string
	Values     object.Object[string]
	Conditions array.Array[condition.Condition]
}

func New() *Set {
	return &Set{}
}

func (get *Set) IsTargetEmpty() bool {
	return get.Target == ""
}

func (get *Set) SetTarget(target string) {
	get.Target = formatter.ToReferenceCase(target)
}
