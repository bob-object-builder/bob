package mariadb

import (
	"salvadorsru/bob/internal/core/drivers/driver"
)

var Driver = driver.Driver{
	Motor:      driver.MariaDB,
	Types:      Types,
	Variations: Variations,
}
