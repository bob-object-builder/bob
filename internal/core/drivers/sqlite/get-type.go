package sqlite

import "salvadorsru/bob/internal/core/drivers"

func GetType(toGet string) string {
	switch drivers.Type(toGet) {
	case drivers.Id, drivers.Int, drivers.Int8, drivers.Int16, drivers.Int32, drivers.Int64:
		return "INTEGER"
	case drivers.Uint, drivers.Uint8, drivers.Uint16, drivers.Uint32, drivers.Uint64:
		return "INTEGER"
	case drivers.Float32, drivers.Float64:
		return "REAL"
	case drivers.String, drivers.String8, drivers.String16, drivers.String32, drivers.String64, drivers.Text, drivers.Date:
		return "TEXT"
	case drivers.Blob:
		return "BLOB"
	case drivers.Time:
		return "INTEGER"
	}

	return ""
}
