package join

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/value"
)

type Selection struct {
	From         string
	Alias        string
	Target       string
	IsExpression bool
}

func NewChild(from string, token value.Token) *Selection {
	s := &Selection{From: from}
	if token.Type == value.Alias {
		s.Alias = token.Value
	} else {
		s.IsExpression = token.Type == value.Expression
		s.Target = token.Value
	}
	return s
}

func (s *Selection) Merge(i any) *failure.Failure { return nil }

func (s *Selection) Parse(token value.Token) *failure.Failure {

	if token.Type == value.Alias {
		s.Alias = token.Value
		return nil
	}

	if s.Target == "" {
		s.IsExpression = token.Type == value.Expression
		s.Target = token.Value
		return nil
	}

	return nil
}
