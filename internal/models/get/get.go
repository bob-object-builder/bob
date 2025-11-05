package get

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/join"
)

const Key = "get"
const GroupKey = "group"

type Get struct {
	Target     string
	Selected   object.Object[string]
	Conditions array.Array[condition.Condition]
	Having     array.Array[condition.Condition]
	Alias      string
	Subqueries array.Array[Get]
	Joins      array.Array[join.Join]
	Groups     array.Array[string]
}

func New(alias ...string) *Get {
	get := &Get{}

	if len(alias) > 0 && alias[0] != "" {
		get.Alias = alias[0]
	}

	return get
}

func (g *Get) HasTarget() bool {
	return g.Target != ""
}

func (g *Get) SetTarget(target string) {
	g.Target = formatter.ToSnakeCase(target)
}

func IsGroup(key string) bool {
	return key == GroupKey
}
