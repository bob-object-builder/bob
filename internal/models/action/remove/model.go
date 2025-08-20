package remove

import (
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
)

type Delete struct {
	Table    string
	Alias    string
	Selected utils.Object[any]
	Filters  []lexer.Instruction
	Join     utils.Object[Join]
}

type Actions []Delete
