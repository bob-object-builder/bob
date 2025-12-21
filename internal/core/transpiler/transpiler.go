package transpiler

import (
	"salvadorsru/bob/internal/core/drivers/driver"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/lib/console"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/drop"
	"salvadorsru/bob/internal/model/get"
	"salvadorsru/bob/internal/model/insert"
	"salvadorsru/bob/internal/model/raw"
	"salvadorsru/bob/internal/model/remove"
	"salvadorsru/bob/internal/model/set"
	"salvadorsru/bob/internal/model/table"
)

type Transpiler struct {
	Driver  driver.Driver
	tables  value.Map[*table.Table]
	actions value.Array[any]
}

func (t *Transpiler) Transpile(query string) (failure *failure.Failure, tables *value.Array[string], actions *value.Array[string]) {
	lx := lexer.Lexer{
		Driver: t.Driver,
	}

	parseError, stack := lx.Parse(query)
	if parseError != nil {
		return parseError, nil, nil
	}

	t.tables = stack.GetTables()
	t.actions = stack.GetActions()

	tablesList := value.NewArray[string]()
	actionsList := value.NewArray[string]()

	for table := range t.tables.Range() {
		tbError, tb := t.TranspileTable(table.Value)
		if tbError != nil {
			console.Log(tbError)
		}

		tablesList.Push(tb)
	}

	for _, action := range t.actions {
		switch v := action.(type) {
		case *get.Get:
			getError, getSQL := TranspileGet(v, false)
			if getError != nil {
				return getError, nil, nil
			}

			actionsList.Push(getSQL)
		case *remove.Delete:
			delError, delSQL := TranspileDelete(v)
			if delError != nil {
				return delError, nil, nil
			}

			actionsList.Push(delSQL)

		case *insert.Insert:
			insertError, insertSQL := TranspileInsert(v)
			if insertError != nil {
				return insertError, nil, nil
			}

			actionsList.Push(insertSQL)

		case *raw.Raw:
			rawSQL := TranspileRaw(v)

			actionsList.Push(rawSQL)

		case *drop.Drop:
			dropSQL := TranspileDrop(v)
			actionsList.Push(dropSQL)

		case *set.Set:
			setError, setSQL := TranspileSet(v)
			if setError != nil {
				return setError, nil, nil
			}

			actionsList.Push(setSQL)

		}
	}

	return nil, tablesList, actionsList
}
