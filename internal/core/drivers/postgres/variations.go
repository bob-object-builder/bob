package postgres

import "salvadorsru/bob/internal/core/drivers/driver"

var Variations = driver.Variations{
	driver.AutoIncrement: "GENERATED ALWAYS AS IDENTITY",
}
