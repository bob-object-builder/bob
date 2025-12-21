package join

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
)

type Join struct {
	ReferencedTable  string
	ReferencedColumn string
	Selected         value.Array[Selection]
	Conditions       condition.Conditions
}

type Joins struct {
	Joins value.Array[Join]
}

func NewJoin() *Join {
	return &Join{ReferencedColumn: "id"}
}

func (js *Joins) Push(j Join) {
	js.Joins.Push(j)
}

func (j *Join) Merge(i any) *failure.Failure {

	switch v := i.(type) {
	case *Selection:
		j.Selected.Push(*v)
	case *condition.Condition:
		j.Conditions.Push(*v)
	}

	return nil
}

func (j *Join) Parse(token string, isParameter bool) {

	if token == kw.NewLine {
		return
	}

	if isParameter {
		if j.ReferencedTable == "" {
			j.ReferencedTable = formatter.ToReferenceCase(token)
			return
		}

		j.ReferencedColumn = formatter.ToReferenceCase(token)
		return
	}
}
