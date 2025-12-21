package remove

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
)

type Delete struct {
	Target     string
	Conditions condition.Conditions
}

func NewDelete() *Delete {
	return &Delete{}
}

func (d *Delete) Parse(token value.Token) {
	if d.Target == "" {
		d.Target = formatter.ToReferenceCase(token.Value)
	}
}

func (d *Delete) Merge(i any) *failure.Failure {
	switch v := i.(type) {
	case *condition.Condition:
		d.Conditions.Push(*v)
		return nil
	}

	return nil
}
