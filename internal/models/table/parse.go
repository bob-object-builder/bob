package table

import (
	"salvadorsru/bob/internal/core/config"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/response"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func parseColumn(table *Table, tableName string, v lexer.Instruction, forcedColumn bool) error {
	columnName := v[0]
	isReference := utils.StartsWithUpper(columnName)

	var columnType string
	var columnAttributes lexer.Instruction

	if len(v) > 1 {
		columnType = v[1]
		columnAttributes = v[2:]
	} else if isReference {
		if columnType == "" {
			columnType = "id"
		}
	}

	isAction := strings.HasPrefix(columnName, config.Use().ActionAlias)
	if isAction {
		switch drivers.Attribute(columnName[1:]) {
		case drivers.Index:
			table.Indexes = append(table.Indexes, Indexable{tableName, v[1:]})
		case drivers.Primary:
			table.PrimaryKeys = append(table.PrimaryKeys, v[1:]...)
		}
		return nil
	}

	isOptional := columnAttributes.Has(string(drivers.Optional))
	if isReference {
		referenceColumnName := utils.PascalToSnakeCase(columnName + "_" + columnType)
		isolated := columnAttributes.Has(string(drivers.Isolated))
		table.References.Set(referenceColumnName, Reference{columnName, columnType, isolated, isOptional})
		columnName = referenceColumnName
	}

	if drivers.Type(columnType) == drivers.Id || drivers.Type(columnType) == drivers.Current {
		tagType, tagAttributes := drivers.UseTag(columnType)
		columnType = tagType
		if !isReference {
			columnAttributes = tagAttributes
		}
	}

	isAutoIncrement := false
	isPrimaryKey := false
	var defaultValue *string

	for i := 0; i < len(columnAttributes); i++ {
		columnAttribute := columnAttributes[i]

		switch drivers.Attribute(columnAttribute) {
		case drivers.Index:
			table.Indexes = append(table.Indexes, NewIndexable(tableName, columnName))
		case drivers.Unique:
			table.Uniques = append(table.Uniques, columnName)
		case drivers.AutoIncrement:
			isAutoIncrement = true
		case drivers.Primary:
			isPrimaryKey = true
			table.PrimaryKeys = append(table.PrimaryKeys, columnName)
		case drivers.Optional:
			isOptional = true
		case drivers.Default:
			if len(columnAttributes) > i+1 {
				firstValue := columnAttributes[i+1]
				if strings.HasPrefix(firstValue, `"`) {
					acc := []string{}

					for p, sliceValue := range columnAttributes[i+1:] {
						acc = append(acc, sliceValue)
						if len(sliceValue) > 0 && sliceValue[len(sliceValue)-1] == '\\' {
							continue
						}

						if strings.HasSuffix(sliceValue, `"`) {
							i += p
							break
						}
					}

					computedValue := strings.Join(acc, " ")
					defaultValue = &computedValue
				} else {
					val := firstValue
					defaultValue = &val
					i += 1
				}
			}
		}
	}

	column := Column{Type: columnType, Attributes: columnAttributes, IsPrimaryKey: isPrimaryKey, IsAutoIncrement: isAutoIncrement, IsOptional: isOptional, Default: defaultValue}

	if forcedColumn {
		table.Columns.Prepend(columnName, &column)
		return nil
	}

	table.Columns.Set(columnName, &column)

	return nil
}

func parseTable(tableId string, block lexer.Block) utils.Object[*Table] {
	tables := utils.Object[*Table]{}
	tableName := utils.PascalToSnakeCase(tableId)
	table := Table{Id: tableId, Name: tableName, Columns: utils.Object[*Column]{}, References: utils.Object[Reference]{}}

	for _, attr := range block.Children() {
		switch v := attr.(type) {
		case lexer.Instruction:
			parseColumn(&table, tableName, v, false)
		case lexer.Block:
			if block.Command == lexer.Table {
				subTableName := v.Actions()[0]
				subTableName = tableName + subTableName
				v.Append(lexer.Instruction{tableId, "id"})
				subTables := parseTable(subTableName, v)
				tables.Merge(subTables.Reverse())
			}
			continue
		}
	}

	if !table.Columns.Has("id") {
		parseColumn(&table, tableId, lexer.Instruction{"id", "id"}, true)
	}

	tables.Set(tableId, &table)

	return tables.Reverse()
}

func makeReference(parsedTables *utils.Object[*Table], referencedTableName string, referencedTableColumn Reference) (error, string, *Column) {
	referencedTable := (*parsedTables).Get(referencedTableColumn.table)
	if referencedTable == nil {
		return response.Error("referenced table not found", referencedTableColumn.table), "", nil
	}
	referecedColumn := referencedTable.Columns.Get(referencedTableColumn.column)
	if referecedColumn == nil {
		return response.Error("referenced column not found", referencedTableColumn.column, "in table", referencedTableColumn.table), "", nil
	}
	columnName := referencedTableName

	return nil, columnName, &Column{
		Type:       referecedColumn.Type,
		Attributes: referecedColumn.Attributes,
	}
}

func Parse(tables lexer.Tables) (error, *utils.Object[*Table]) {
	parsedTables := utils.Object[*Table]{}

	for _, tableName := range tables.Order {
		block, _ := tables.Get(tableName)
		tables := parseTable(tableName, block)
		parsedTables.Merge(tables)
	}

	for _, parsedTableName := range parsedTables.Order {
		parsedTable := parsedTables.Get(parsedTableName)

		for _, referencedTableName := range parsedTable.References.Order {
			referencedTableColumn := parsedTable.References.Get(referencedTableName)
			columnsError, columnName, newColumn := makeReference(&parsedTables, referencedTableName, referencedTableColumn)
			if columnsError != nil {
				return columnsError, nil
			}

			newColumn.IsOptional = referencedTableColumn.optional

			parsedTable.Columns.Set(columnName, newColumn)
		}
	}

	return nil, &parsedTables
}
