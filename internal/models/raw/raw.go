package raw

import "salvadorsru/bob/internal/lib/value/array"

const Key = "raw"

type Raw struct {
	Lines array.Array[string]
}

func New() *Raw {
	return &Raw{}
}
