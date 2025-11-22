package driver

import (
	"salvadorsru/bob/internal/core/driver/mariadb"
	"salvadorsru/bob/internal/core/driver/postgres"
	"salvadorsru/bob/internal/core/driver/sqlite"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/models/function"
	"salvadorsru/bob/internal/models/table"
)

type Motor string

const (
	SQLite   Motor = "sqlite"
	MariaDB  Motor = "mariadb"
	Postgres Motor = "postgres"
	MySQL    Motor = "mysql"
)

type Driver struct {
	Motor Motor
}

func GetDriver(motor string) (*failure.Failure, *Driver) {
	switch Motor(motor) {
	case SQLite:
		return nil, &Driver{SQLite}
	case MariaDB:
		return nil, &Driver{MariaDB}
	case Postgres:
		return nil, &Driver{Postgres}
	case MySQL:
		return nil, &Driver{MySQL}
	default:
		return failure.UnknownDriver(motor), nil
	}
}

func (d *Driver) GetType(token table.Type) (*failure.Failure, table.Type) {
	switch d.Motor {
	case SQLite:
		return sqlite.Types.GetType(token)
	case MariaDB, MySQL:
		return mariadb.Types.GetType(token)
	case Postgres:
		return postgres.Types.GetType(token)
	}

	panic("unselected driver")
}

func (d *Driver) GetCaller(token string) (*failure.Failure, function.Caller) {
	switch d.Motor {
	case MariaDB, MySQL:
		return mariadb.Callers.GetCaller(token)
	}

	return failure.UndefinedCaller, ""
}
