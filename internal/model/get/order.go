package get

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"strings"
)

type Order struct {
	From       string
	Direction  string
	Target     string
	HasNulls   bool
	NullsFirst bool
}

func NewOrder(from string, direction string) *Order {
	return &Order{
		From:      from,
		Direction: direction,
	}
}

func (o *Order) Merge(i any) *failure.Failure {
	return nil
}

func (o *Order) Parse(token string) *failure.Failure {
	if o.Target == "" {
		if token == kw.NewLine {
			return failure.UndefinedOrderTarget
		}

		o.Target = strings.ToLower(token)
		return nil
	}

	if token == kw.NewLine {
		return nil
	}

	if !o.HasNulls {

		if token == kw.Nulls {
			o.HasNulls = true
			return nil
		} else {
			return failure.UndefinedNullsDefinition
		}

	} else {
		switch token {
		case kw.NullsFirst:
			o.NullsFirst = true
			return nil
		case kw.NullsLast:
			o.NullsFirst = false
			return nil
		default:
			return failure.UndefinedNullsDefinition
		}

	}

}
