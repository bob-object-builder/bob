package postgres

import "salvadorsru/bob/internal/core/drivers/driver"

var Driver = driver.Driver{
	Motor:      driver.Postgres,
	Types:      Types,
	Variations: Variations,
}
