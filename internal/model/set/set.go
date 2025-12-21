package set

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/condition"
)

type Set struct {
	Target     string
	Values     value.Map[string]
	Conditions condition.Conditions
	acc        value.Array[string]
}

func NewSet() *Set {
	return &Set{
		Values:     *value.NewMap[string](),
		Conditions: condition.Conditions{},
	}
}

func (s *Set) Parse(token value.Token) *failure.Failure {
	if s.Target == "" {
		s.Target = token.Value
		return nil
	}

	if token.Value == kw.NewLine {
		if s.acc.Length() == 1 {
			return failure.UndefinedValueOnSetter(s.Target, s.acc[0])
		}

		if s.acc.Length() == 2 {
			key := s.acc[0]
			value := s.acc[1]
			s.Values.Set(key, value)
			s.acc.Clean()
		}
		return nil
	}

	s.acc.Push(token.Value)
	return nil
}

func (d *Set) Merge(i any) *failure.Failure {
	switch v := i.(type) {
	case *condition.Condition:
		d.Conditions.Push(*v)
		return nil
	}

	return nil
}
