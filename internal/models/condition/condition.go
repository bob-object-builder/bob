package condition

import "salvadorsru/bob/internal/lib/value/array"

type ConditionKey string
type Comparator string

const (
	If ConditionKey = "if"
	Or ConditionKey = "or"
)

const (
	And        Comparator = "&&"
	Else       Comparator = "||"
	Like       Comparator = "like"
	Equal      Comparator = "="
	BiggerThan Comparator = ">"
	LowerThan  Comparator = "<"
	Different  Comparator = "!="
)

type Condition struct {
	Condition  ConditionKey
	Table      string
	Target     string
	Comparator Comparator
	And        array.Array[string]
	Else       array.Array[string]
}

func IsCondition(s string) bool {
	switch s {
	case string(If), string(Or):
		return true
	}
	return false
}

func IsComparator(s string) bool {
	switch Comparator(s) {
	case And, Else, Like, Equal, BiggerThan, LowerThan, Different:
		return true
	}
	return false
}

func IsAnd(s string) bool {
	return s == string(And)
}

func IsElse(s string) bool {
	return s == string(Else)
}
