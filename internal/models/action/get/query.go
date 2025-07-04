package get

import (
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func NewQuery(block lexer.Block) Get {
	get := Get{Filters: []lexer.Instruction{}}
	tableName := block.Actions()[0]
	get.Table = utils.PascalToSnakeCase(tableName)

	for _, child := range block.Children() {
		switch v := child.(type) {
		case lexer.Instruction:
			firstItem := v[0]
			isAlias := strings.HasSuffix(v[0], ":")

			if isAlias {
				alias := firstItem[0 : len(firstItem)-1]

				if len(v) == 1 {
					get.Alias = alias
				} else {
					get.Selected.Set(alias, v[1:])
				}
				continue
			}

			isOperatorion := drivers.Operator(firstItem) == drivers.If || drivers.Operator(firstItem) == drivers.Or || drivers.Operator(firstItem) == drivers.Group
			if isOperatorion {
				get.Filters = append(get.Filters, v)
				continue
			}

			get.Selected.Set(firstItem, v[1:])
		case lexer.Block:

			if v.ActionIs(lexer.LeftJoin, lexer.LeftJoinAlias) {
				var joinColumn string
				if len(v.Actions()) > 1 {
					joinColumn = v.Actions()[1]
				} else {
					joinColumn = "id"
				}

				joinTableName := v.Actions()[0]
				joinQuery := NewQuery(v)

				get.Join.Set(joinTableName, Join{
					Direction: lexer.Left,
					Table:     utils.PascalToSnakeCase(joinTableName),
					On:        joinColumn,
					Query:     &joinQuery,
				})

			}

			subQuery := NewQuery(v)
			get.Selected.Set(subQuery.Alias, subQuery)
		}
	}

	return get
}
