package postgresql

import "salvadorsru/bob/internal/core/drivers"

func GetType(toGet string) string {
	switch drivers.Type(toGet) {
	case drivers.Int:
		return "INTEGER"
	case drivers.Int8:
		return "SMALLINT"
	case drivers.Int16:
		return "SMALLINT"
	case drivers.Int32:
		return "INTEGER"
	case drivers.Int64:
		return "BIGINT"

	case drivers.Uint:
		return "INTEGER"
	case drivers.Uint8:
		return "SMALLINT"
	case drivers.Uint16:
		return "INTEGER"
	case drivers.Uint32:
		return "BIGINT"
	case drivers.Uint64:
		return "NUMERIC"

	case drivers.Float32:
		return "REAL"
	case drivers.Float64:
		return "DOUBLE PRECISION"

	case drivers.String:
		return "VARCHAR(255)"
	case drivers.String8:
		return "VARCHAR(8)"
	case drivers.String16:
		return "VARCHAR(16)"
	case drivers.String32:
		return "VARCHAR(32)"
	case drivers.String64:
		return "VARCHAR(64)"

	case drivers.Text:
		return "TEXT"
	case drivers.Blob:
		return "BYTEA"
	case drivers.Date:
		return "DATE"
	case drivers.Time:
		return "TIME"

	case drivers.Id:
		return "SERIAL"
	case drivers.Boolean:
		return "BOOLEAN"

	default:
		return ""
	}
}
