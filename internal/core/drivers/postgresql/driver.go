package postgresql

import (
	"salvadorsru/bob/internal/core/drivers"
)

var Driver = drivers.Driver{
	Motor:        drivers.PostgreSQL,
	GetType:      GetType,
	GetAttribute: GetAttribute,
	GetFunction:  GetFunction,
	GetLiteral:   GetLiteral,
}
