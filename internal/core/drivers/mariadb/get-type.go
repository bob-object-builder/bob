package mariadb

import "salvadorsru/bob/internal/core/drivers"

func GetType(toGet string) string {
	switch drivers.Type(toGet) {
	case drivers.Id, drivers.Int:
		return "INT"
	case drivers.Int8:
		return "TINYINT"
	case drivers.Int16:
		return "SMALLINT"
	case drivers.Int32:
		return "INT"
	case drivers.Int64:
		return "BIGINT"

	case drivers.Uint:
		return "INT UNSIGNED"
	case drivers.Uint8:
		return "TINYINT UNSIGNED"
	case drivers.Uint16:
		return "SMALLINT UNSIGNED"
	case drivers.Uint32:
		return "INT UNSIGNED"
	case drivers.Uint64:
		return "BIGINT UNSIGNED"

	case drivers.Float32:
		return "FLOAT"
	case drivers.Float64:
		return "DOUBLE"

	case drivers.String, drivers.String8, drivers.String16, drivers.String32, drivers.String64:
		switch drivers.Type(toGet) {
		case drivers.String8:
			return "VARCHAR(8)"
		case drivers.String16:
			return "VARCHAR(16)"
		case drivers.String32:
			return "VARCHAR(32)"
		case drivers.String64:
			return "VARCHAR(64)"
		default:
			return "VARCHAR(255)"
		}
	case drivers.Text:
		return "TEXT"
	case drivers.Blob:
		return "BLOB"
	case drivers.Date:
		return "DATE"
	case drivers.Time:
		return "TIME"
	case drivers.Boolean:
		return "BOOLEAN"
	default:
		return ""
	}
}
