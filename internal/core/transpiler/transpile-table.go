package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/literal"
	"salvadorsru/bob/internal/models/table"
	"strings"
)

func TranspileIndex(table, column string) string {
	return fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s);", table, column, table, column)
}

func (t Transpiler) TranspileReference(ref table.Reference) (*failure.Failure, *table.Column, string) {
	referencedTable := ref.Table
	referencedColumn := ref.Column
	columnName := fmt.Sprintf("%s_%s", referencedTable, referencedColumn)
	foreignKey := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", columnName, referencedTable, referencedColumn)

	tb := t.Tables.Get(referencedTable)
	if tb == nil {
		return failure.UndefinedReferenceTable(referencedTable), nil, ""
	}

	col := tb.Columns.Get(referencedColumn)
	if col == nil {
		return failure.UndefinedReferencedColumn(referencedColumn, referencedTable), nil, ""
	}

	column := table.Column{
		Name:     columnName,
		Type:     col.Type,
		Optional: ref.Optional,
		Default:  ref.Default,
	}

	onUpdate := array.New[string]()
	if ref.OnDeleteCascade {
		onUpdate.Push("ON DELETE CASCADE")
	}

	if ref.OnUpdateCascade {
		onUpdate.Push("ON UPDATE CASCADE")
	}

	if onUpdate.Length() > 0 {
		foreignKey += "\n" + formatter.IndentLines(strings.Join(*onUpdate, "\n"), 2)
	}

	return nil, &column, foreignKey
}

func (t *Transpiler) TranspileColumn(col table.Column) (*failure.Failure, string) {
	if col.Type == "" {
		return failure.UndefinedTypeForColumn(col.GetName()), ""
	}

	typeError, dbType := t.GetType(col.Type)
	if typeError != nil {
		return typeError, ""
	}

	query := fmt.Sprintf("%s %s", col.GetName(), dbType)

	if col.Primary {
		query += " PRIMARY KEY"
	}

	if col.AutoIncrement {
		switch t.SelectedDriver {
		case "sqlite":
			query += " AUTOINCREMENT"
		case "mysql":
			query += " AUTO_INCREMENT"
		}
	}

	if !col.Optional {
		query += " NOT NULL"
	}

	if col.Default != "" {
		query += fmt.Sprintf(" DEFAULT %s", literal.GetLiteral(col.Default))
	}

	return nil, query
}

func (t *Transpiler) TranspileTable(tb table.Table) (*failure.Failure, array.Array[string]) {
	var (
		columns = array.New[string]()
		indexes = array.New[string]()
		uniques = array.New[string]()
		output  = array.New[string]()
	)

	tableName := tb.GetName()

	for column := range tb.Columns.Range() {
		col := column.Value
		colName := col.GetName()

		if col.Index {
			indexes.Push(TranspileIndex(tableName, colName))
		}

		if col.Unique {
			uniques.Push(colName)
		}

		columnError, colSQL := t.TranspileColumn(col)
		if columnError != nil {
			return columnError, nil
		}

		columns.Push(formatter.Indent(colSQL))
	}

	for _, ref := range tb.References {
		err, colRef, fkDef := t.TranspileReference(ref)
		if err != nil {
			return err, nil
		}

		err, colSQL := t.TranspileColumn(*colRef)
		if err != nil {
			return err, nil
		}

		columns.Push(formatter.Indent(colSQL))
		columns.Push(formatter.Indent(fkDef))
	}

	if uniques.Length() > 0 {
		columns.Push(formatter.Indent(fmt.Sprintf("UNIQUE(%s)", strings.Join(*uniques, ", "))))
	}

	output.Push(fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s (\n%s\n);",
		tableName,
		strings.Join(*columns, ",\n"),
	))

	if indexes.Length() > 0 {
		output.Push(*indexes...)
	}

	return nil, *output
}
