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
const LimitKey = "limit"
const OffsetKey = "offset"
const EveryField = "*"
const SpreadEveryField = "..."

type Get struct {
	Target     string
	Selected   object.Object[string]
	Conditions array.Array[condition.Condition]
	Having     array.Array[condition.Condition]
	Alias      string
	Subqueries array.Array[Get]
	Joins      array.Array[join.Join]
	Groups     array.Array[string]
	Limit      string
	Offset     string
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

func IsLimit(key string) bool {
	return key == LimitKey
}

func IsOffset(key string) bool {
	return key == OffsetKey
}

func IsEveryField(selected string) bool {
	return selected == EveryField || selected == SpreadEveryField
}
