package mariadb

import (
	"salvadorsru/bob/internal/core/drivers"
)

func GetAttribute(attr string) string {
	switch drivers.Attribute(attr) {
	case drivers.AutoIncrement:
		return "AUTO_INCREMENT"
	case drivers.Required:
		return "NOT NULL"
	}
	return ""
}
