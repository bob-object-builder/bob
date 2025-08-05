package set

import (
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
)

type Set struct {
	Table   string
	Alias   string
	Values  utils.Object[any]
	Filters []lexer.Instruction
	Join    utils.Object[Join]
}

type Actions []Set
