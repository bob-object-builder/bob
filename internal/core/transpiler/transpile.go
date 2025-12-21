package transpiler

import (
	"salvadorsru/bob/internal/core/drivers/driver"
	"salvadorsru/bob/internal/core/drivers/mariadb"
	"salvadorsru/bob/internal/core/drivers/postgres"
	"salvadorsru/bob/internal/core/drivers/sqlite"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/value"
)

func Transpile(driverName string, query string) (*failure.Failure, value.Array[string], value.Array[string]) {
	driverError, driver := GetDriver(driverName)
	if driverError != nil {
		return driverError, nil, nil
	}

	t := &Transpiler{
		Driver: *driver,
	}

	error, tables, actions := t.Transpile(query)
	if error != nil {
		return error, nil, nil
	}

	return nil, *tables, *actions
}

func GetDriver(driverName string) (*failure.Failure, *driver.Driver) {
	switch driver.Motor(driverName) {
	case driver.SQLite:
		return nil, &sqlite.Driver
	case driver.MariaDB:
		return nil, &mariadb.Driver
	case driver.MySQL:
		return nil, &mariadb.Driver
	case driver.Postgres:
		return nil, &postgres.Driver
	default:
		return failure.UnknownDriver(driverName), nil
	}
}
