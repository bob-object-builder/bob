package postgresql

import "salvadorsru/bob/internal/core/drivers"

func GetType(toGet string) string {
	switch drivers.Type(toGet) {
	case drivers.Int:
		return "INTEGER"
	case drivers.Int8:
		return "SMALLINT" // no tiene tinyint, usamos SMALLINT
	case drivers.Int16:
		return "SMALLINT"
	case drivers.Int32:
		return "INTEGER"
	case drivers.Int64:
		return "BIGINT"

	case drivers.Uint:
		return "INTEGER" // no unsigned en PG, hay que controlar en app
	case drivers.Uint8:
		return "SMALLINT"
	case drivers.Uint16:
		return "INTEGER"
	case drivers.Uint32:
		return "BIGINT"
	case drivers.Uint64:
		return "NUMERIC" // para uint64 sin signo, usar NUMERIC

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
		// Serial for auto increment (in tables use SERIAL or IDENTITY)
		return "SERIAL"
	case drivers.CreatedAt:
		return "TIMESTAMP DEFAULT CURRENT_TIMESTAMP"

	default:
		return string(toGet)
	}
}
