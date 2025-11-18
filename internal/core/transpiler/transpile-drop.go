package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/models/drop"
)

func (t Transpiler) TranspileDrop(d drop.Drop) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;", d.Target)
}
