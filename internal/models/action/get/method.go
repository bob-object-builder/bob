package get

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

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

			subJoins, subJoinQueries := ParseJoins(driver, toJoin.Query, selected)
			joins = append(joins, subJoins...)
			joinQueries = append(joinQueries, subJoinQueries...)
		}
	}

	return joins, joinQueries
}

func ParseOperations(actions ...*Get) []string {
	operations := []string{}
	isGroup := false
	isFirstOperationInsideGroup := false

	for _, a := range actions {
		for i, filter := range a.Filters {
			sentence := []string{}
			operator := filter[0]
			subject := filter[1]

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
							sentence = append(sentence, stringAcc)
							onString = false
						}

						continue
					}

					if utils.IsStringStart(token) {
						stringAcc = token
						onString = true

						if utils.IsStringEnd(token) {
							sentence = append(sentence, stringAcc)
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
