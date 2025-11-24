package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/set"
	"strings"
)

func (t Transpiler) TranspileSet(s set.Set) (*failure.Failure, string) {
	query := "UPDATE \n%s\nSET\n%s%s;"

	conditionString := ""
	var conditionError *failure.Failure
	conditionError, conditionString = t.TranspileConditions(s.Conditions, false)
	if conditionError != nil {
		return conditionError, ""
	}

	setsList := array.New[string]()
	for set := range s.Values.Range() {
		setsList.Push(fmt.Sprintf("%s = %s", set.Key, set.Value))
	}

	sets := formatter.IndentLines(strings.Join(*setsList, ",\n"))

	return nil, fmt.Sprintf(query,
		formatter.Indent(s.Target),
		sets,
		conditionString,
	)
}
