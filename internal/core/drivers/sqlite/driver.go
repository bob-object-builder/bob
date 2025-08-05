package sqlite

import (
	"salvadorsru/bob/internal/core/drivers"
)

var Driver = drivers.Driver{
	Motor:        drivers.SQLite,
	GetType:      GetType,
	GetAttribute: GetAttribute,
	GetFunction:  GetFunction,
	GetLiteral:   GetLiteral,
}
