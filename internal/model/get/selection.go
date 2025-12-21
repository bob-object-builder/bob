package get

import (
	"salvadorsru/bob/internal/lib/value"
)

type Selection struct {
	From         string
	Selected     value.Token
	IsExpression bool
}

func (s *Selection) isGetChildValue() {}

func NewSelection(tk value.Token) *Selection {
	return &Selection{Selected: tk}
}
