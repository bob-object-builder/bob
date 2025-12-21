package table

import (
	"salvadorsru/bob/internal/core/kw"
)

type Reference struct {
	Target     string
	Column     string
	IsOptional bool
}

func NewReference(target string) *Reference {
	return &Reference{
		Target: target,
	}
}

func (r *Reference) Parse(token string) {
	if token == kw.NewLine {
		if r.Column == "" {
			r.Column = "id"
		}
		return
	}

	if r.Target == "" {
		r.Target = token
		return
	}

	if r.Column == "" {
		r.Column = token
		return
	}

	if token == "optional" {
		r.IsOptional = true
		return
	}
}
