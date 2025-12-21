package sqlite

import (
	"salvadorsru/bob/internal/core/drivers/driver"
)

var Driver = driver.Driver{
	Motor:      driver.SQLite,
	Types:      Types,
	Variations: Variations,
}
