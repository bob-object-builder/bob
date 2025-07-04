package postgresql

import "salvadorsru/bob/internal/core/drivers"

func GetFunction(fn string) string {
	switch drivers.Function(fn) {
	case drivers.CONCAT:
		return "CONCAT"
	}
	return ""
}
