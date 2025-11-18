package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/models/join"
)

func (t Transpiler) TranspileLeftJoin(tableName string, join join.Join) string {
	query := "LEFT JOIN %s ON %s = %s.%s"
	tableField := fmt.Sprintf("%s.%s_%s", tableName, join.Target, join.On)
	return fmt.Sprintf(query, join.Target, tableField, join.Target, join.On)
}
