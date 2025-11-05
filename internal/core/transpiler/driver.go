package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers/mariadb"
	"salvadorsru/bob/internal/core/drivers/postgre"
	"salvadorsru/bob/internal/core/drivers/sqlite"
	"salvadorsru/bob/internal/models/table"
)

type Driver string

const (
	SQLite  Driver = "sqlite"
	MariaDB Driver = "mariadb"
	Postgre Driver = "postgres"
	MySQL   Driver = "mysql"
)

func GetDriver(driver string) (error, Driver) {
	switch Driver(driver) {
	case SQLite:
		return nil, SQLite
	case MariaDB:
		return nil, MariaDB
	case Postgre:
		return nil, Postgre
	case MySQL:
		return nil, MySQL
	default:
		return fmt.Errorf("unknown driver: %s", driver), ""
	}
}

func (t *Transpiler) SetDriver(driver Driver) {
	t.SelectedDriver = driver
}

func (t *Transpiler) GetType(token table.Type) (table.Type, error) {
	switch t.SelectedDriver {
	case SQLite:
		return sqlite.Types.GetType(token)
	case MariaDB, MySQL:
		return mariadb.Types.GetType(token)
	case Postgre:
		return postgre.Types.GetType(token)
	}

	panic("unselected driver")
}
