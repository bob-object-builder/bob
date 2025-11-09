package function

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"strings"
)

type FunctionKey string
type ContextKey string

const (
	avgKey    FunctionKey = "avg"
	sumKey    FunctionKey = "sum"
	minKey    FunctionKey = "min"
	maxKey    FunctionKey = "max"
	countKey  FunctionKey = "count"
	concatKey FunctionKey = "concat" // MySQL: GROUP_CONCAT, PostgreSQL: string_agg, SQLite: group_concat
)

const (
	openKey  ContextKey = "("
	closeKey ContextKey = ")"
)

type Function struct {
	Key          FunctionKey
	Args         []string
	Subfunctions array.Array[Function]
}

func IsFunction(s string) bool {
	s = strings.TrimSuffix(s, "(")
	s = strings.ToLower(s)
	for _, f := range []string{
		string(openKey),
		string(avgKey), string(sumKey), string(minKey),
		string(maxKey), string(countKey), string(concatKey),
	} {
		if strings.HasPrefix(s, f) {
			return true
		}
	}
	return false
}

func IsFunctionStart(token string) bool {
	return strings.Contains(token, "(") && IsFunction(token)
}

func ReconstructFunction(tokens []string, startIndex int) (string, int) {
	var (
		tokenBuffer = array.New[string]()
		openParens  int
	)

	for i := startIndex; i < len(tokens); i++ {
		t := tokens[i]
		tokenBuffer.Push(t)

		openParens += strings.Count(t, "(")
		openParens -= strings.Count(t, ")")

		if openParens <= 0 {
			return strings.Join(*tokenBuffer, " "), i
		}
	}

	return strings.Join(*tokenBuffer, " "), len(tokens) - 1
}

func PrefixParameters(prefix string, target string) string {
	return formatter.PrefixWith(
		prefix,
		target,
		[]string{
			string(avgKey),
			string(sumKey),
			string(minKey),
			string(maxKey),
			string(countKey),
			string(concatKey),
		},
	)
}
