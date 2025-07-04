package get

import "salvadorsru/bob/internal/core/lexer"

type Join struct {
	Direction lexer.Direction
	Table     string
	On        string
	Query     *Get
}
