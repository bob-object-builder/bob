package driver

import (
	"salvadorsru/bob/internal/core/failure"
)

type Variation string

const (
	Length        Variation = "length"
	AutoIncrement Variation = "auto_increment"
)

type Variations map[Variation]string

func (v *Variations) Get(token string) string {
	if val, ok := (*v)[Variation(token)]; ok {
		return val
	}

	return token
}

func (v *Variations) MustGet(token string) (*failure.Failure, string) {
	val, ok := (*v)[Variation(token)]

	if !ok {
		return failure.UndefinedKeyword("property", token), ""
	}

	return nil, val
}
