package drop

import "salvadorsru/bob/internal/lib/formatter"

const Key = "drop"

type Drop struct {
	Target string
}

func New() *Drop {
	return &Drop{}
}

func (g *Drop) HasTarget() bool {
	return g.Target != ""
}

func (g *Drop) SetTarget(target string) {
	g.Target = formatter.ToSnakeCase(target)
}
