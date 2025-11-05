package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/table"
	"strings"
)

func TranspileIndex(table, column string) string {
	return fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s);", table, column, table, column)
}

func (t Transpiler) TranspileReference(ref table.Reference) (string, string, error) {
	referencedTable := ref.Table
	referencedColumn := ref.Column

	tb := t.Tables.Get(referencedTable)
	if tb == nil {
		return "", "", fmt.Errorf("undefined reference table '%s'", referencedTable)
	}

	col := tb.Columns.Get(referencedColumn)
	if col == nil {
		return "", "", fmt.Errorf("undefined referenced column: '%s' in table '%s'", referencedColumn, referencedTable)
	}

	dbType, err := t.GetType(col.Type)
	if err != nil {
		return "", "", err
	}

	columnName := fmt.Sprintf("%s_%s", referencedTable, referencedColumn)
	columnDef := fmt.Sprintf("%s %s", columnName, dbType)
	if !ref.IsOptional {
		columnDef += " NOT NULL"
	}

	foreignKey := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", columnName, referencedTable, referencedColumn)
	return columnDef, foreignKey, nil
}

func (t *Transpiler) TranspileColumn(col table.Column) (string, error) {
	if col.Type == "" {
		return "", fmt.Errorf("undefined type for column: %s", col.GetName())
	}

	dbType, err := t.GetType(col.Type)
	if err != nil {
		return "", err
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
		query += fmt.Sprintf(" DEFAULT %s", col.Default)
	}

	return query, nil
}

func (t *Transpiler) TranspileTable(tb table.Table) (error, array.Array[string]) {
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

		colSQL, err := t.TranspileColumn(col)
		if err != nil {
			return err, nil
		}

		columns.Push(formatter.Indent(colSQL))
	}

	for _, ref := range tb.References {
		colDef, fkDef, err := t.TranspileReference(ref)
		if err != nil {
			return err, nil
		}

		columns.Push(formatter.Indent(colDef))
		columns.Push(formatter.Indent(fkDef))
	}

	if uniques.Length() > 0 {
		columns.Push(formatter.Indent(fmt.Sprintf("UNIQUE(%s)", strings.Join(*uniques, ", "))))
	}

	output.Push(fmt.Sprintf(
		"CREATE TABLE %s (\n%s\n);",
		tableName,
		strings.Join(*columns, ",\n"),
	))

	if indexes.Length() > 0 {
		output.Push(*indexes...)
	}

	return nil, *output
}
