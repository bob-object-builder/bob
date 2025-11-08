package transpiler

import (
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"salvadorsru/bob/internal/models/get"
	"salvadorsru/bob/internal/models/insert"
	"salvadorsru/bob/internal/models/raw"
	"salvadorsru/bob/internal/models/remove"
	"salvadorsru/bob/internal/models/table"
)

type Transpiler struct {
	Tables         object.Object[table.Table]
	Actions        array.Array[any]
	SelectedDriver Driver
}

type transpileMode string

func (t Transpiler) Transpile() (error, *TranspiledTable, *TranspiledActions) {
	tablesError, tables := t.TranspileTables()
	if tablesError != nil {
		return tablesError, nil, nil
	}

	actionsError, actions := t.TranspileActions()
	if actionsError != nil {
		return actionsError, nil, nil
	}

	return nil, tables, actions
}

func (t Transpiler) TranspileTables() (error, *TranspiledTable) {
	tables := TranspiledTable{}

	for table := range t.Tables.Range() {
		tableError, table := t.TranspileTable(table.Value)
		if tableError != nil {
			return tableError, nil
		}
		tables.Push(table...)
	}

	return nil, &tables
}

func (t Transpiler) TranspileActions() (error, *TranspiledActions) {
	actions := TranspiledActions{}

	for _, action := range t.Actions {
		switch a := action.(type) {
		case get.Get:
			actions.Push(t.TranspileGet(a))
		case insert.Insert:
			actions.Push(t.TranspileInsert(a))
		case remove.Remove:
			actions.Push(t.TranspileRemove(a))
		case raw.Raw:
			actions.Push(t.TranspileRaw(a))
		}
	}

	return nil, &actions
}
