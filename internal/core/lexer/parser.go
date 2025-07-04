package lexer

import (
	"strings"
)

func Parser(query string) Program {
	lines := strings.SplitSeq(query, "\n")

	pile := NewPile()
	program := NewProgram()
	currentInstruction := NewInstructions()
	atActions := true
	blockDepth := 0

line:
	for line := range lines {

		if line == "" {
			continue
		}

		for token := range strings.FieldsSeq(line) {
			if token[0] == '#' {
				continue line
			}

			if IsCommand(token) {
				pile.Add(token)
				atActions = true
				continue
			}

			if token == "{" {
				blockDepth++
				atActions = false
				continue
			}

			if atActions {
				pile.GetLast().Append(token)
				continue
			}

			if token == "}" {
				blockDepth--

				if blockDepth == 0 {
					program.Append(*pile.GetLast())
					continue
				}

				pile.MergeLast()
				continue
			}

			currentInstruction.Append(token)
		}

		if len(currentInstruction) != 0 {
			pile.GetLast().Append(currentInstruction)
			currentInstruction = NewInstructions()
		}
	}

	return *program
}
