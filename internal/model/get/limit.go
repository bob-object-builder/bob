package get

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/checker"
)

type Limit struct {
	Limit  string
	Offset string
}

func NewLimit() *Limit {
	return &Limit{}
}

func (l *Limit) Merge(i any) *failure.Failure {
	return nil
}

func (l *Limit) Parse(token string) *failure.Failure {
	if l.Limit == "" {
		if !checker.IsNumber(token) {
			return failure.LimitMustBeNumeric
		}
		l.Limit = token
		return nil
	}

	if token == kw.NewLine {
		return nil
	}

	if l.Offset == "" {
		if !checker.IsNumber(token) {
			return failure.LimitMustBeNumeric
		}
		l.Offset = token
		return nil
	}

	return nil
}
