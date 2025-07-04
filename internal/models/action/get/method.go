package get

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"strconv"
	"strings"
)

func isValue(input string) bool {
	trimmed := strings.TrimSpace(input)

	if len(trimmed) == 0 {
		return false
	}

	if (strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") && len(trimmed) > 1) ||
		(strings.HasPrefix(trimmed, "'") && strings.HasSuffix(trimmed, "'") && len(trimmed) > 1) {
		return true
	}

	if _, err := strconv.Atoi(trimmed); err == nil {
		return true
	}

	if _, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return true
	}

	lower := strings.ToLower(trimmed)
	if lower == "true" || lower == "false" {
		return true
	}

	if lower == "null" {
		return true
	}

	return false
}

func ParseSelections(driver drivers.Driver, a *Get, selected *[]string, isJoin bool) {
	for _, selectedName := range a.Selected.Keys() {
		value := a.Selected.Get(selectedName)
		switch v := value.(type) {
		case lexer.Instruction:
			if len(v) == 1 && !drivers.HasFunction(v[0]) {
				selectedName = fmt.Sprintf("%s.%s as %s", v[0], a.Table, selectedName)
			} else {
				computedSelected := strings.Join(v, " ")
				computedSelected = drivers.ReplaceFunction(computedSelected, driver.GetFunction)
				if computedSelected == "" {

					if drivers.HasFunction(selectedName) {
						selectedName = fmt.Sprintf("%s", selectedName)
					} else {
						if isJoin {
							tableName := a.Table
							if len(tableName) > 2 && strings.HasSuffix(tableName, "ies") {
								tableName = tableName[:len(tableName)-3] + "y"
							} else if len(tableName) > 0 && tableName[len(tableName)-1] == 's' {
								tableName = tableName[:len(tableName)-1]
							}
							selectedName = fmt.Sprintf("%s.%s as %s_%s", a.Table, selectedName, tableName, selectedName)
						} else {
							selectedName = fmt.Sprintf("%s.%s as %s", a.Table, selectedName, selectedName)
						}
					}

				} else {
					selectedName = fmt.Sprintf("%s as %s_%s", computedSelected, a.Table, selectedName)
				}
			}

			selectedName = utils.Indent(selectedName)
			*selected = append(*selected, selectedName)
		case Get:
			query := utils.IndentLines(v.ToQuery(driver))
			selectedName = fmt.Sprintf("(\n%s\n) as %s", query, selectedName)
			selectedName = utils.IndentLines(selectedName)
			*selected = append(*selected, selectedName)
			continue
		}
	}
}

func ParseJoins(driver drivers.Driver, a *Get, selected *[]string) ([]string, []*Get) {
	joins := []string{}
	joinQueries := []*Get{}

	for _, toJoinName := range a.Join.Keys() {
		toJoin := a.Join.Get(toJoinName)

		switch toJoin.Direction {
		case lexer.Left:
			leftOn := fmt.Sprintf("%s.%s_%s", a.Table, toJoin.Table, toJoin.On)
			rightOn := fmt.Sprintf("%s.%s", toJoin.Table, toJoin.On)

			joinSentence := fmt.Sprintf("LEFT JOIN %s ON %s = %s", toJoin.Table, leftOn, rightOn)
			joins = append(joins, joinSentence)

			ParseSelections(driver, toJoin.Query, selected, true)
			joinQueries = append(joinQueries, toJoin.Query)

			// RECURSIVAMENTE PARSEAMOS LOS JOINS INTERNOS
			subJoins, subQueries := ParseJoins(driver, toJoin.Query, selected)
			joins = append(joins, subJoins...)
			joinQueries = append(joinQueries, subQueries...)
		}
	}

	return joins, joinQueries
}

func ParseOperations(actions ...*Get) []string {
	operations := []string{}

	addOp := func(newOp ...string) {
		operations = append(operations, strings.Join(newOp, " "))
	}

	addLnOp := func(newOp ...string) {
		operations = append(operations, "\n"+strings.Join(newOp, " "))
	}

	isFirst := true
	isGroup := false

	for _, a := range actions {
		for _, filter := range a.Filters {
			if len(operations) == 0 {
				addLnOp("WHERE")
			}
			var targetColumn string
			var latestOperation string
			for _, token := range filter {
				op := drivers.Operator(token)

				switch op {
				case drivers.If:
					if isFirst {
						if isGroup {
							addLnOp("HAVING")
						}
						isFirst = false
					} else {
						addLnOp("AND")
					}
					if len(filter) > 1 {
						targetColumn = filter[1]
					}
					continue
				case drivers.Or:
					addLnOp("OR")
					if len(filter) > 1 {
						targetColumn = filter[1]
					}
				case drivers.Group:
					addLnOp("GROUP BY")
					isFirst = true
					isGroup = true
					continue
				case drivers.And:
					addLnOp("AND", fmt.Sprintf("%s.%s", a.Table, targetColumn), latestOperation)
				case drivers.Else:
					addLnOp("OR", fmt.Sprintf("%s.%s", a.Table, targetColumn), latestOperation)
				default:
					if isValue(token) {
						addOp(token)
					} else if drivers.IsOperator(token) {
						addOp(token)
						latestOperation = token
					} else {
						if utils.StartsWithUpper(token) {
							addOp(utils.PascalToSnakeCase(token))
						} else {
							addOp(fmt.Sprintf("%s.%s", a.Table, token))
						}
					}
				}
			}
		}
	}

	return operations
}

func (a Get) ToQuery(driver drivers.Driver) string {
	query := "SELECT\n%s \nFROM %s"

	selected := []string{}

	ParseSelections(driver, &a, &selected, false)

	if len(selected) < 1 {
		selected = append(selected, "*")
	}

	joins, joinQueries := ParseJoins(driver, &a, &selected)

	query = fmt.Sprintf(query, strings.Join(selected, ",\n"), a.Table)

	operations := ParseOperations(append(joinQueries, &a)...)

	if len(joins) > 0 {
		joinsSentence := strings.Join(joins, "\n")
		query += "\n" + joinsSentence
	}

	if len(operations) > 0 {
		operationSentence := strings.Join(operations, " ")
		query += " " + operationSentence
	}

	return query
}
