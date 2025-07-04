package table

import (
	"fmt"
	"salvadorsru/bob/internal/core/utils"
)

type Reference struct {
	table           string
	column          string
	onDeleteCascade bool
}

func (r *Reference) toQuery() string {
	sql := "FOREIGN KEY (%s) REFERENCES %s(%s)"
	if r.onDeleteCascade {
		sql += " ON DELETE CASCADE"
	}
	tableName := utils.PascalToSnakeCase(r.table)
	columnName := fmt.Sprintf("%s_%s", tableName, r.column)
	return fmt.Sprintf(sql, columnName, tableName, r.column)
}
