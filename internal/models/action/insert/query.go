package insert

import (
	"salvadorsru/bob/internal/core/console"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func NewQuery(block lexer.Block) New {
	new := New{}
	tableName := block.Actions()[0]
	new.Table = utils.PascalToSnakeCase(tableName)

	bulkInsertion := len(block.Actions()) > 1

	formatColumnName := func(columnName string) string {
		if utils.StartsWithUpper(columnName) {
			columnName = strings.ReplaceAll(columnName, ".", "_")
			columnName = utils.PascalToSnakeCase(columnName)
		}

		return columnName
	}

	if bulkInsertion {
		for _, columnName := range block.Actions()[1:] {
			columnName = formatColumnName(columnName)
			new.Fields = append(new.Fields, columnName)
		}

		for _, child := range block.Children() {
			values := []string{}

			if v, ok := child.(lexer.Instruction); ok {
				columnString := []string{}

				for i, slice := range v {

					if i > len(block.Actions()[1:]) {
						break
					}

					if utils.IsStringStart(slice) {
						if utils.IsStringEnd(slice) {
							columnString = []string{}
							columnString = append(columnString, slice)
							values = append(values, utils.FormatQuote(strings.Join(columnString, " ")))
						} else {
							columnString = append(columnString, slice)
						}
						continue
					}

					if utils.IsStringEnd(slice) {
						columnString = append(columnString, slice)
						values = append(values, utils.FormatQuote(strings.Join(columnString, " ")))
						continue
					}

					values = append(values, slice)

				}
			}

			new.Values = append(new.Values, values)
		}
	} else {
		values := []string{}
		fieldSet := utils.NewObject[int]()

		for position, child := range block.Children() {
			if v, ok := child.(lexer.Instruction); ok {
				columnName := formatColumnName(v[0])
				columnValue := v[1:]
				columnName = strings.TrimSuffix(columnName, ":")

				if fieldSet.Has(columnName) {
					console.Log(columnName, columnName)
					values[fieldSet.Get(columnName)] = utils.FormatQuote(strings.Join(columnValue, " "))
					continue
				}

				fieldSet.Set(columnName, position)
				new.Fields = append(new.Fields, columnName)
				values = append(values, utils.FormatQuote(strings.Join(columnValue, " ")))
			}
		}

		new.Values = append(new.Values, values)
	}

	return new
}
