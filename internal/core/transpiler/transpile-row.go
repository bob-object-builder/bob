package transpiler

import (
	"salvadorsru/bob/internal/models/raw"
	"strings"
)

func (t Transpiler) TranspileRaw(r raw.Raw) string {
	return strings.Join(r.Lines, " ")
}
