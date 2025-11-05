package transpiler

import (
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"salvadorsru/bob/internal/models/get"
	"salvadorsru/bob/internal/models/insert"
	"salvadorsru/bob/internal/models/remove"
	"salvadorsru/bob/internal/models/table"
	"strings"
)

type Transpiler struct {
	Tables         object.Object[table.Table]
	Actions        array.Array[any]
	SelectedDriver Driver
}

type transpileMode string

func (t Transpiler) Transpile() (error, string, string) {
	tablesError, tables := t.TranspileTables()
	if tablesError != nil {
		return tablesError, "", ""
	}

	actionsError, actions := t.TranspileActions()
	if actionsError != nil {
		return actionsError, "", ""
	}

	return nil, tables, actions
}

func (t Transpiler) TranspileTables() (error, string) {
	tables := array.New[string]()

	for table := range t.Tables.Range() {
		tableError, table := t.TranspileTable(table.Value)
		if tableError != nil {
			return tableError, ""
		}
		tables.Push(table...)
	}

	return nil, strings.Join(*tables, "\n\n")
}

func (t Transpiler) TranspileActions() (error, string) {
	actions := array.New[string]()

	for _, action := range t.Actions {
		switch a := action.(type) {
		case get.Get:
			actions.Push(t.TranspileGet(a))
		case insert.Insert:
			actions.Push(t.TranspileInsert(a))
		case remove.Remove:
			actions.Push(t.TranspileRemove(a))
		}
	}

	return nil, strings.Join(*actions, "\n\n")
}
