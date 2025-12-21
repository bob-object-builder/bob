package raw

import "salvadorsru/bob/internal/lib/value"

type Raw struct {
	Query value.Array[string]
}

func NewRaw() *Raw {
	return &Raw{}
}

func (r *Raw) Parse(token value.Token) {
	r.Query.Push(token.Value)
}
