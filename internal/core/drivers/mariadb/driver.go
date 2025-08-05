package mariadb

import (
	"salvadorsru/bob/internal/core/drivers"
)

var Driver = drivers.Driver{
	Motor:        drivers.MariaDB,
	GetType:      GetType,
	GetAttribute: GetAttribute,
	GetFunction:  GetFunction,
	GetLiteral:   GetLiteral,
}
