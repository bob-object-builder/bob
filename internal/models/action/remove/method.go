package remove

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/response"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func ParseOperations(actions ...*Delete) []string {
	operations := []string{}
	isFirstOperation := true

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

			comparator := filter[2]
			onString := false
			stringAcc := ""

			if i == 0 || isFirstOperation {
				sentence = append(sentence, "\nWHERE")

				isFirstOperation = false
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
						sentence = append(sentence, "\nOR")
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

			operations = append(operations, strings.Join(sentence, " "))
		}
	}

	return operations
}

func (a Delete) ToQuery(driver drivers.Driver) (error, string) {
	query := "DELETE FROM %s"

	query = fmt.Sprintf(query, a.Table)

	if len(a.Join.Keys()) > 0 {
		return response.Error("you cannot use joins at delete"), ""
	}

	deletes := []*Delete{&a}
	operations := ParseOperations(deletes...)

	if len(operations) > 0 {
		operationSentence := strings.Join(operations, " ")
		query += " " + operationSentence
	}

	return nil, query
}
