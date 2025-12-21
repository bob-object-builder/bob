package transpiler

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/model/condition"
)

func TranspileGroup(g *condition.Group) string {
	return "\nGROUP BY\n" + formatter.IndentLines(g.Targets.Join(",\n"))
}
