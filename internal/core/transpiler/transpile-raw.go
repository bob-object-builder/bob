package transpiler

import (
	"salvadorsru/bob/internal/model/raw"
	"strings"
)

func TranspileRaw(raw *raw.Raw) string {
	return strings.TrimSpace(raw.Query.Join(" "))
}
