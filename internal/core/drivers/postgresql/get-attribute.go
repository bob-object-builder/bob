package postgresql

import (
	"salvadorsru/bob/internal/core/drivers"
)

func GetAttribute(attr string) string {
	switch drivers.Attribute(attr) {
	case drivers.Optional:
		return "NOT NULL"
	}
	return ""
}
