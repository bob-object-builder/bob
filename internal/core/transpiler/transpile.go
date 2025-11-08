package transpiler

import (
	"salvadorsru/bob/internal/core/lexer"
)

func Transpile(driver Driver, query string) (error, *TranspiledTable, *TranspiledActions) {
	l := lexer.New()
	parseError, tables, actions := l.Parse(query)

	if parseError != nil {
		return parseError, nil, nil
	}

	t := Transpiler{
		Tables:         *tables,
		Actions:        *actions,
		SelectedDriver: driver,
	}

	trampileError, transpiledTables, transpiledActions := t.Transpile()
	if trampileError != nil {
		return trampileError, nil, nil
	}

	return nil, transpiledTables, transpiledActions
}
