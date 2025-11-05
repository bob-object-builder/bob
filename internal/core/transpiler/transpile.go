package transpiler

import (
	"salvadorsru/bob/internal/core/lexer"
)

func Transpile(driver Driver, query string) (error, string, string) {
	l := lexer.New()
	parseError, tables, actions := l.Parse(query)

	if parseError != nil {
		return parseError, "", ""
	}

	t := Transpiler{
		Tables:         *tables,
		Actions:        *actions,
		SelectedDriver: driver,
	}

	trampileError, transpiledTables, transpiledActions := t.Transpile()
	if trampileError != nil {
		return trampileError, "", ""
	}

	return nil, transpiledTables, transpiledActions
}
