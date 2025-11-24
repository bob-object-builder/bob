package join

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/order"
)

const (
	LeftJoinKey = "->"
)

type Join struct {
	Target     string
	Direction  string
	From       string
	To         string
	On         string
	Selected   object.Object[string]
	Subjoins   array.Array[Join]
	Conditions array.Array[condition.Condition]
	Having     array.Array[condition.Condition]
	Groups     array.Array[string]
	Capturing  bool
	Orders     array.Array[order.Order]
}

func IsJoin(join string) bool {
	return join == LeftJoinKey
}

func NewLeftJoin() *Join {
	return &Join{Capturing: false}
}

func (j *Join) IsTargetEmpty() bool {
	return j.Target == ""
}

func (j *Join) IsOnEmpty() bool {
	return j.On == ""
}

func (j *Join) SetTarget(target string) {
	j.Target = formatter.ToReferenceCase(target)
}

func (j *Join) SetOn(on ...string) {
	if len(on) == 0 {
		j.On = "id"
		return
	}

	j.On = formatter.ToReferenceCase(on[0])
}
