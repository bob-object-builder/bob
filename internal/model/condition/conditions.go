package condition

import (
	"salvadorsru/bob/internal/lib/value"
)

type Conditions struct {
	Conditions value.Array[Condition]
}

func (cs *Conditions) Push(condition Condition) {
	cs.Conditions.Push(condition)
}

func (cs *Conditions) Prepend(condition Condition) {
	cs.Conditions.Prepend(condition)
}
