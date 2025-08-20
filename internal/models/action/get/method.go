package get

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/response"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func ParseSelections(driver drivers.Driver, a *Get, selected *[]string, isJoin bool) error {
	for _, selectedName := range a.Selected.Keys() {
		value := a.Selected.Get(selectedName)

		switch v := value.(type) {
		case lexer.Instruction:
			if len(v) == 1 && !drivers.HasFunction(v[0]) {
				selectedName = fmt.Sprintf("%s.%s as %s", a.Table, v[0], selectedName)
			} else if selectedName == "*" || selectedName == "..." {
				selectedName = a.Table + ".*"
			} else {
				computedSelected := strings.Join(v, " ")
				computedSelected = drivers.ReplaceFunction(computedSelected, driver.GetFunction)
				if computedSelected == "" {

					if !drivers.HasFunction(selectedName) {
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
			subQueryError, subQuery := v.ToQuery(driver)
			if subQueryError != nil {
				return subQueryError
			}

			query := utils.IndentLines(subQuery)

			if selectedName == "" {
				return response.Error("you must set an alias in subquery")
			}

			selectedName = fmt.Sprintf("(\n%s\n) as %s", query, selectedName)
			selectedName = utils.IndentLines(selectedName)
			*selected = append(*selected, selectedName)
			continue
		}
	}

	return nil
}

func ParseJoins(driver drivers.Driver, a *Get, selected *[]string) (error, []string, []*Get) {
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

			parseSelectionError := ParseSelections(driver, toJoin.Query, selected, true)
			if parseSelectionError != nil {
				return parseSelectionError, nil, nil
			}
			joinQueries = append(joinQueries, toJoin.Query)

			subJoinError, subJoins, subJoinQueries := ParseJoins(driver, toJoin.Query, selected)
			if subJoinError != nil {
				return subJoinError, nil, nil
			}
			joins = append(joins, subJoins...)
			joinQueries = append(joinQueries, subJoinQueries...)
		}
	}

	return nil, joins, joinQueries
}

func ParseOperations(actions ...*Get) []string {
	operations := []string{}
	isGroup := false
	isFirstOperationInsideGroup := false

	if len(actions) == 0 {
		return operations
	}

	for _, a := range actions {
		for i, filter := range a.Filters {
			if len(filter) < 2 {
				// optionally log or skip
				continue
			}

			sentence := []string{}
			operator := filter[0]
			subject := a.Table + "." + filter[1]

			if operator == string(drivers.Group) {
				sentence = append(sentence, "\nGROUP BY", subject)
				isGroup = true
				isFirstOperationInsideGroup = true
			} else {

				comparator := filter[2]
				onString := false
				stringAcc := ""

				if i == 0 || isFirstOperationInsideGroup {
					if isGroup {
						sentence = append(sentence, "\nHAVING")
					} else {
						sentence = append(sentence, "\nWHERE")
					}

					isFirstOperationInsideGroup = false
				} else {
					if operator == string(drivers.Or) {
						sentence = append(sentence, "\nOR")
					} else {
						sentence = append(sentence, "\nAND")
					}
				}

				for _, token := range filter[2:] {

					if onString {
						stringAcc += " " + token

						if utils.IsStringEnd(token) {
							sentence = append(sentence, utils.FormatQuote(stringAcc))
							onString = false
						}

						continue
					}

					if utils.IsStringStart(token) {
						stringAcc = token
						onString = true

						if utils.IsStringEnd(token) {
							sentence = append(sentence, utils.FormatQuote(stringAcc))
							onString = false
						}

						continue
					}

					if drivers.IsOperator(token) {
						op := drivers.Operator(token)

						switch op {
						case drivers.Equal:
						case drivers.Or:
						case drivers.Else:
							sentence = append(sentence, "OR")
						case drivers.And:
							sentence = append(sentence, "AND")
						}

						sentence = append(sentence, subject, comparator)

						continue
					}

					if utils.IsValue(token) {
						sentence = append(sentence, token)
						continue
					}
				}
			}

			operations = append(operations, strings.Join(sentence, " "))
		}
	}

	return operations
}

func (a Get) ToQuery(driver drivers.Driver) (error, string) {
	query := "SELECT\n%s \nFROM %s"

	selected := []string{}

	parseSelections := ParseSelections(driver, &a, &selected, false)
	if parseSelections != nil {
		return parseSelections, ""
	}

	if len(selected) < 1 {
		selected = append(selected, "*")
	}

	parseJoinsError, joins, joinQueries := ParseJoins(driver, &a, &selected)
	if parseJoinsError != nil {
		return parseJoinsError, ""
	}

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

	return nil, query
}
