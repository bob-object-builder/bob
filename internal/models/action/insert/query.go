package insert

import (
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func NewQuery(block lexer.Block) New {
	new := New{}
	tableName := block.Actions()[0]
	new.Table = utils.PascalToSnakeCase(tableName)

	bulkInsertion := len(block.Actions()) > 1

	if bulkInsertion {
		for _, columnName := range block.Actions()[1:] {
			if utils.StartsWithUpper(columnName) {
				columnName = strings.ReplaceAll(columnName, ".", "_")
				columnName = utils.PascalToSnakeCase(columnName)
			}

			new.Fields = append(new.Fields, columnName)
		}

		inString := false
		for _, child := range block.Children() {
			values := []string{}

			if v, ok := child.(lexer.Instruction); ok {
				columnString := []string{}

				for _, slice := range v {
					if inString {
						columnString = append(columnString, slice)

						isEscaped := slice[len(slice)-1] == '\\'
						if !isEscaped && strings.HasSuffix(slice, `"`) {
							inString = false
							values = append(values, strings.Join(columnString, " "))
						}

					} else {
						inString = strings.HasPrefix(slice, `"`)

						if inString {
							columnString = append(columnString, slice)

							if strings.HasSuffix(slice, `"`) {
								inString = false
								values = append(values, strings.Join(columnString, " "))
							}
						} else {
							values = append(values, slice)
						}
					}

				}
			}

			new.Values = append(new.Values, values)
		}
	} else {
		values := []string{}
		fieldSet := utils.NewObject[int]()

		for position, child := range block.Children() {
			if v, ok := child.(lexer.Instruction); ok {
				columnName := utils.PascalToSnakeCase(v[0])
				columnValue := v[1:]
				columnName = strings.TrimSuffix(columnName, ":")

				if fieldSet.Has(columnName) {
					values[fieldSet.Get(columnName)] = strings.Join(columnValue, " ")
					continue
				}

				fieldSet.Set(columnName, position)
				new.Fields = append(new.Fields, columnName)
				values = append(values, strings.Join(columnValue, " "))
			}
		}

		new.Values = append(new.Values, values)
	}

	return new
}
