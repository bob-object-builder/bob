package drivers

import (
	"fmt"
	"salvadorsru/bob/internal/core/console"
	"strings"
)

type Function string

const (
	AVG    Function = "avg"
	SUM    Function = "sum"
	MIN    Function = "min"
	MAX    Function = "max"
	COUNT  Function = "count"
	CONCAT Function = "concat" // MySQL: GROUP_CONCAT, PostgreSQL: string_agg, SQLite: group_concat
)

func ReplaceFunction(s string, replacer func(toReplace string) string) string {
	if replacer == nil {
		console.Panic("ReplaceFunction: replacer function is nil")
		return s
	}
	functions := []Function{AVG, SUM, MIN, MAX, COUNT, CONCAT}
	for _, fn := range functions {
		functionName := string(fn)
		replacement := replacer(functionName)
		if replacement != "" {
			s = strings.ReplaceAll(s, functionName, replacement)
		}
	}

	if s != "" {
		return fmt.Sprintf("(%s)", s)
	}

	return ""
}

func HasFunction(str string) bool {
	return strings.Contains(str, "(")
}
