package sqlite

import "salvadorsru/bob/internal/core/drivers"

func GetFunction(string string) string {
	switch drivers.Function(string) {
	case drivers.CONCAT:
		return "GROUP_CONCAT"
	}

	return ""
}
