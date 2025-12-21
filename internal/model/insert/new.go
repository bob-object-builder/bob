package insert

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"strings"
)

type Insert struct {
	onDefinition bool
	IsBulk       bool
	Target       string
	Fields       value.Array[string]
	Values       value.Array[value.Array[string]]
	acc          value.Array[string]
}

func NewInsert() *Insert {
	return &Insert{
		onDefinition: true,
	}
}

func (n *Insert) Parse(token value.Token, isParameter bool) *failure.Failure {

	if isParameter {
		if n.Target == "" {
			n.Target = strings.ToLower(token.Value)
		} else {
			n.Fields.Push(
				formatter.ToReferenceCase(token.Value),
			)
		}
		return nil
	}

	if n.onDefinition {
		n.IsBulk = n.Fields.Length() != 0
		n.onDefinition = false
		return nil
	}

	if token.Value == kw.NewLine {
		if n.acc.Length() > 0 {
			if n.IsBulk {
				n.Values.Push(n.acc)
			} else {
				field := n.acc[0]
				n.Fields.Push(formatter.ToReferenceCase(field))

				ref := n.Values.GetLast()
				if ref == nil {
					n.Values.Push(n.acc[1:])
				} else {
					ref.Push(n.acc[1:]...)
				}
			}

			n.acc.Clean()
		}

		return nil
	}

	n.acc.Push(
		token.Value,
	)

	return nil
}
