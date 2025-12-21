package get

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/value"
	"strings"
)

type GetChildValue interface {
	isGetChildValue()
}

type Child struct {
	From  string
	Alias string
	Value GetChildValue
}

func NewGetChild(from string, alias string) *Child {
	return &Child{
		From:  from,
		Alias: strings.TrimSuffix(alias, kw.Alias),
	}
}

func (g *Child) Parse(token value.Token) {
	g.Value = NewSelection(token)
}

func (g *Child) Merge(i any) *failure.Failure {

	switch v := i.(type) {
	case *Get:
		g.Value = v
	}

	return nil
}
