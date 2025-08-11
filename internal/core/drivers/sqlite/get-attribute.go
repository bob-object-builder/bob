package sqlite

import (
	"salvadorsru/bob/internal/core/drivers"
)

func GetAttribute(string string) string {
	switch drivers.Attribute(string) {
	case drivers.AutoIncrement:
		return "AUTOINCREMENT"
	case drivers.Optional:
		return "NOT NULL"
	}
	return ""
}
