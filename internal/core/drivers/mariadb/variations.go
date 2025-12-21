package mariadb

import (
	"salvadorsru/bob/internal/core/drivers/driver"
)

var Variations = driver.Variations{
	driver.Length:        "CHAR_LENGTH",
	driver.AutoIncrement: "AUTO_INCREMENT",
}
