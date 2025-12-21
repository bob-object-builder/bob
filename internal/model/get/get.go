package get

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
	"salvadorsru/bob/internal/model/join"
)

type Get struct {
	Target     string
	Children   value.Array[Child]
	Group      condition.Group
	HasGroup   bool
	Conditions condition.Conditions
	Havings    condition.Conditions
	Joins      join.Joins
	Orders     value.Array[Order]
}

func NewGet() *Get {
	return &Get{}
}

func (p *Get) isGetChildValue() {}

func (g *Get) Merge(i any) *failure.Failure {

	switch v := i.(type) {
	case *Child:
		g.Children.Push(*v)
		return nil
	case *join.Join:
		g.Joins.Push(*v)
		return nil
	case *Order:
		g.Orders.Push(*v)
		return nil
	case *condition.Condition:
		if g.HasGroup {
			g.Havings.Push(*v)
		} else {
			g.Conditions.Push(*v)
		}
		return nil
	case *condition.Group:
		g.Group = *v
		g.HasGroup = true
		return nil
	}

	return nil
}

func (g *Get) Parse(token string) {
	if g.Target == "" {
		g.Target = formatter.ToReferenceCase(token)
		return
	}
}
