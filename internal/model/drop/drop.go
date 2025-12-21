package drop

import (
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
)

type Drop struct {
	Target string
}

func NewDrop() *Drop {
	return &Drop{}
}

func (d *Drop) Parse(token value.Token) {
	if d.Target == "" {
		d.Target = formatter.ToReferenceCase(token.Value)
	}
}
