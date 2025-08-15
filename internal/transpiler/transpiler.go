package transpiler

import (
	"salvadorsru/bob/internal/core/drivers"
	mariadb "salvadorsru/bob/internal/core/drivers/mariadb"
	postgresql "salvadorsru/bob/internal/core/drivers/postgresql"
	sqlite "salvadorsru/bob/internal/core/drivers/sqlite"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/response"
	"salvadorsru/bob/internal/models/action"
	"salvadorsru/bob/internal/models/table"
)

func getDriver(motor drivers.Motor) (error, *drivers.Driver) {
	var driver drivers.Driver
	switch motor {
	case drivers.SQLite:
		driver = sqlite.Driver
	case drivers.MariaDB:
		driver = mariadb.Driver
	case drivers.PostgreSQL:
		driver = postgresql.Driver
	default:
		return response.Error("Unknown driver '%s'. Supported drivers: mariadb, postgresql, sqlite", motor), nil
	}

	return nil, &driver
}

type TranspilerMaker func(motor drivers.Motor, input string) (error, string, string)

func makeTranspiler(transpileTables bool, transpileActions bool) TranspilerMaker {
	return func(motor drivers.Motor, input string) (error, string, string) {
		var (
			tablesQueries = ""
			actionQueries = ""
		)

		var tablesQueriesError, actionQueriesError error

		driverError, driver := getDriver(motor)
		if driverError != nil {
			return driverError, "", ""
		}

		program := lexer.Parser(input)

		if transpileTables {
			tablesQueriesError, tablesQueries = table.Transpile(*driver, program.Tables)
		}

		if transpileActions {
			actionQueriesError, actionQueries = action.Transpile(*driver, program.Actions)
		}

		if tablesQueriesError != nil {
			return tablesQueriesError, "", ""
		}

		if actionQueriesError != nil {
			return actionQueriesError, "", ""
		}

		return nil, tablesQueries, actionQueries
	}
}

var Transpile = makeTranspiler(true, true)
var TranspileTables = makeTranspiler(true, false)
var TranspileActions = makeTranspiler(false, true)
