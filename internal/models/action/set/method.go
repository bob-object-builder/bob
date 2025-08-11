package set

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func ParseSelections(driver drivers.Driver, a *Set, selected *[]string, isJoin bool) {
	for _, selectedName := range a.Values.Keys() {
		value := a.Values.Get(selectedName)

		switch v := value.(type) {
		case lexer.Instruction:
			if len(v) == 0 {
				continue
			}

			if !drivers.HasFunction(v[0]) {
				selectedName = fmt.Sprintf("%s = %s", selectedName, utils.FormatQuote(strings.Join(v, " ")))
			}

			selectedName = utils.Indent(selectedName)
			*selected = append(*selected, selectedName)
		case Set:
			query := utils.IndentLines(v.ToQuery(driver))
			selectedName = fmt.Sprintf("(\n%s\n) as %s", query, selectedName)
			selectedName = utils.IndentLines(selectedName)
			*selected = append(*selected, selectedName)
			continue
		}
	}
}

func ParseJoins(driver drivers.Driver, a *Set, selected *[]string) ([]string, []*Set) {
	joins := []string{}
	joinQueries := []*Set{}

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

func ParseOperations(actions ...*Set) []string {
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

func (a Set) ToQuery(driver drivers.Driver) string {
	var query string

	if driver.Motor == drivers.PostgreSQL {
		query = "UPDATE\n%s\nSET\n%s %s %s"
	} else {
		query = "UPDATE\n%s %s \nSET \n%s %s"
	}

	selected := []string{}

	ParseSelections(driver, &a, &selected, false)

	if len(selected) < 1 {
		selected = append(selected, utils.Indent("*"))
	}

	joins, joinQueries := ParseJoins(driver, &a, &selected)

	operations := ParseOperations(append(joinQueries, &a)...)

	var joinsSentence string = ""
	if driver.Motor != drivers.SQLite {
		if len(joins) > 0 {
			joinsSentence = "\n" + strings.Join(joins, "\n")
		}
	}

	var operationSentence string
	if len(operations) > 0 {
		operationSentence = strings.Join(operations, " ")
	}

	if driver.Motor == drivers.MariaDB || driver.Motor == drivers.SQLite {
		query = fmt.Sprintf(query, utils.Indent(a.Table), joinsSentence, strings.Join(selected, ",\n"), operationSentence)
	} else {
		query = fmt.Sprintf(query, utils.Indent(a.Table), strings.Join(selected, ",\n"), joinsSentence, operationSentence)
	}

	return query
}
