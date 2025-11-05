package lexer

import (
	"fmt"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/function"
	"strings"
)

func (l *Lexer) ParseReferences(target string) string {
	tokens := l.tokens

	var (
		currentList   = array.New[string]()
		parsingString bool
		buffer        = array.New[string]()
	)

	prefix := func(str string) string {
		return strings.ToLower(function.PrefixParameters(fmt.Sprintf("%s.", target), str))
	}

	push := func(str string) {
		currentList.Push(prefix(str))
	}

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch {
		case parsingString || checker.IsStringStart(token):
			parsingString = true
			buffer.Push(token)
			if checker.IsStringEnd(token) {
				parsingString = false
				fullString := strings.Join(*buffer, " ")
				buffer.Clean()
				push(formatter.String(fullString))
			}

		case function.IsFunctionStart(token):
			completeFunction, lastIndex := function.ReconstructFunction(tokens, i)
			i = lastIndex
			push(completeFunction)

		default:
			push(token)
		}
	}

	return strings.Join(*currentList, " ")
}
