package condition

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
)

type Condition struct {
	From       string
	Target     string
	Condition  string
	Comparator Comparator
	ValueList  *value.Array[string]
	And        value.Array[string]
	Else       value.Array[string]
}

func NewCondition(from string, condition string) *Condition {
	return &Condition{
		From:      formatter.ToReferenceCase(from),
		Condition: condition,
	}
}

func (c *Condition) Merge() {}

func (c *Condition) Parse(token value.Token) *failure.Failure {
	v := token.Value

	if c.Target == "" {
		if token.Type == value.Key {
			c.Target = formatter.ToReferenceCase(v)
		} else {
			c.Target = v
		}
		return nil
	}

	if c.Comparator == "" {
		if !IsComparator(v) {
			return failure.UndefinedConditionComparator(v)
		}

		c.Comparator = Comparator(v)
		return nil
	}

	switch v {
	case kw.AndSymbol:
		c.ValueList = &c.And
		return nil
	case kw.OrSymbol:
		c.ValueList = &c.Else
		return nil
	}

	if c.ValueList == nil {
		c.ValueList = &c.And
	}

	if token.Type == value.Key {
		c.ValueList.Push(formatter.ToReferenceCase(v))
	} else {
		c.ValueList.Push(v)
	}

	return nil
}
