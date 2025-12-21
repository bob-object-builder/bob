package condition

import "salvadorsru/bob/internal/lib/value"

type Group struct {
	Targets value.Array[string]
}

func NewGroup(target string) *Group {
	return &Group{}
}

func (g *Group) Parse(token string) {
	g.Targets.Push(token)
}
