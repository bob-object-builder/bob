package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers/driver"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
	"salvadorsru/bob/internal/model/table"
)

func (t *Transpiler) TranspileTable(tb *table.Table) (*failure.Failure, string) {
	if tb.Name == "" {
		return failure.TableNameIsEmpty, ""
	}

	sql := "CREATE TABLE %s (\n%s\n);"

	columns := value.NewArray[string]()
	references := value.NewArray[string]()
	indexes := value.NewArray[string]()
	unique := value.NewArray[string]()
	primaries := value.NewArray[string]()

	for v := range tb.References.Range() {
		r := v.Value
		column := r.Column
		if column == "" {
			column = "id"
		}
		columnsName := fmt.Sprintf("%s_%s", r.Target, column)
		ref := fmt.Sprintf("FOREIGN KEY (%s) REFERENCES %s(%s)", columnsName, r.Target, column)
		ref = formatter.Indent(ref)

		referecedTb := *t.tables.Get(r.Target)
		if referecedTb == nil {
			return failure.InvalidTableColumn, ""
		}

		referencedColumn := referecedTb.Columns.Get(column)
		if referencedColumn == nil {
			return failure.InvalidTableColumn, ""
		}

		referencedColumn.Name = columnsName
		referencedColumn.IsAutoIncrement = false
		referencedColumn.IsPrimary = false
		referencedColumn.IsReference = true
		referencedColumn.IsOptional = r.IsOptional

		tb.Columns.Set(columnsName, *referencedColumn)
		references.Push(ref)
	}

	for v := range tb.Columns.Range() {
		c := v.Value

		tp := t.Driver.Types.Get(c.Type)
		if tp == "" {
			return failure.TypeNotFound(c.Type), ""
		}

		if c.Type == string(driver.Id) && !c.IsReference {
			c.IsAutoIncrement = true
			c.IsPrimary = true
		}

		if c.Type == string(driver.Current) {
			c.IsCurrent = true
		}

		col := fmt.Sprintf("%s %s", c.Name, tp)

		if t.Driver.Motor == driver.SQLite {
			if c.IsPrimary {
				col += " PRIMARY KEY"
			}
		} else {
			if c.IsPrimary {
				primaries.Push(c.Name)
			}
		}

		if c.IsAutoIncrement {

			autoIncrementError, autoIncrement := t.Driver.Variations.MustGet(string(driver.AutoIncrement))

			if autoIncrementError != nil {
				return autoIncrementError, ""
			}

			if autoIncrement != "" {
				col += " " + autoIncrement
			}

		}

		if c.IsCurrent {
			c.Default = "CURRENT_TIMESTAMP"
		}

		if c.IsUnique {
			unique.Push(c.Name)
		}

		if !c.IsOptional {
			col += " NOT NULL"
		}

		if c.Default != "" {
			col += fmt.Sprintf(" DEFAULT %s", c.Default)
		}

		if c.IsIndex {
			indexes.Push(
				fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s);", tb.Name, c.Name, tb.Name, c.Name),
			)
		}

		if t.Driver.Motor == driver.SQLite {
			var limit int
			switch driver.Type(c.Type) {
			case driver.String8:
				limit = 8
			case driver.String16:
				limit = 16
			case driver.String32:
				limit = 32
			case driver.String64:
				limit = 64
			case driver.Int8:
				tb.Checks.Prepend(condition.Condition{
					And: *value.NewArray(fmt.Sprintf("%s >= -128 AND %s <= 127", c.Name, c.Name)),
				})
			case driver.Int16:
				tb.Checks.Prepend(condition.Condition{
					And: *value.NewArray(fmt.Sprintf("%s >= -32768 AND %s <= 32767", c.Name, c.Name)),
				})
			case driver.Int32:
				tb.Checks.Prepend(condition.Condition{
					And: *value.NewArray(fmt.Sprintf("%s >= -2147483648 AND %s <= 2147483647", c.Name, c.Name)),
				})
			case driver.Int64:
				tb.Checks.Prepend(condition.Condition{
					And: *value.NewArray(fmt.Sprintf("%s >= -9223372036854775808 AND %s <= 9223372036854775807", c.Name, c.Name)),
				})
			}

			if limit > 0 {
				tb.Checks.Prepend(condition.Condition{
					And: *value.NewArray(fmt.Sprintf("LENGTH(%s) <= %d", c.Name, limit)),
				})
			}
		}

		col = formatter.Indent(col)
		columns.Push(col)
	}

	checksError, checksString := TranspileChecks(&tb.Checks, false, nil)
	if checksError != nil {
		return checksError, ""
	}

	columns.Push(*references...)

	if checksString != "" {
		columns.Push(checksString)
	}

	if primaries.Length() > 0 {
		primaryKey := fmt.Sprintf("PRIMARY KEY (%s)", primaries.Join(", "))
		primaryKey = formatter.Indent(primaryKey)
		columns.Push(primaryKey)
	}

	output := columns.Join(",\n")
	output = fmt.Sprintf(sql, tb.Name, output)

	if indexes.Length() > 0 {
		output += "\n\n" + indexes.Join("\n")
	}

	if unique.Length() > 0 {
		columns.Push(
			fmt.Sprintf("UNIQUE(%s)", unique.Join(", ")),
		)
	}

	return nil, output
}
