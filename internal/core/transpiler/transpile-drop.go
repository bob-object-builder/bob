package transpiler

import (
	"salvadorsru/bob/internal/model/drop"
)

func TranspileDrop(d *drop.Drop) string {
	return "DROP TABLE IF EXISTS " + d.Target + ";"
}
