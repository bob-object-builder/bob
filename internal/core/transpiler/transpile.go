package transpiler

import (
	"salvadorsru/bob/internal/core/driver"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/lexer"
)

func Transpile(motor string, query string) (*failure.Failure, *TranspiledTable, *TranspiledActions) {
	driverError, driver := driver.GetDriver(motor)
	if driverError != nil {
		return driverError, nil, nil
	}

	l := lexer.New(*driver)
	parseError, tables, actions := l.Parse(query)

	if parseError != nil {
		return parseError, nil, nil
	}

	t := Transpiler{
		Tables:  *tables,
		Actions: *actions,
		Driver:  *driver,
	}

	trampileError, transpiledTables, transpiledActions := t.Transpile()
	if trampileError != nil {
		return trampileError, nil, nil
	}

	return nil, transpiledTables, transpiledActions
}
