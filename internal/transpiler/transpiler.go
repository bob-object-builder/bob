package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	mariadb "salvadorsru/bob/internal/core/drivers/mariadb"
	postgresql "salvadorsru/bob/internal/core/drivers/postgresql"
	sqlite "salvadorsru/bob/internal/core/drivers/sqlite"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/models/action"
	"salvadorsru/bob/internal/models/table"
)

func Transpile(motor drivers.Motor, input string) string {
	var driver drivers.Driver
	switch motor {
	case drivers.SQLite:
		driver = sqlite.Driver
	case drivers.MariaDB:
		driver = mariadb.Driver
	case drivers.PostgreSQL:
		driver = postgresql.Driver
	default:
		fmt.Printf("Error: Unknown driver '%s'. Supported drivers: mariadb, postgresql, sqlite\n", motor)
		return ""
	}
	program := lexer.Parser(input)
	parsedTables, tablesQueries := table.Transpile(driver, program.Tables)
	getQueries := action.Transpile(driver, parsedTables, program.Actions)

	return fmt.Sprintf("%s\n\n%s\n", tablesQueries, getQueries)
}
