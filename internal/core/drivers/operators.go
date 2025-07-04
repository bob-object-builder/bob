package drivers

type Operator string

const (
	And        Operator = "&&"
	Else       Operator = "||"
	Like       Operator = "like"
	Equal      Operator = "="
	BiggerThan Operator = ">"
	LowerThan  Operator = "<"
	Different  Operator = "!="
	If         Operator = "if"
	Or         Operator = "or"
	As         Operator = "as"
	Group      Operator = "group"
)

func IsOperator(s string) bool {
	switch Operator(s) {
	case And, Else, Like, Equal, BiggerThan, LowerThan, Different, If, Or, As, Group:
		return true
	default:
		return false
	}
}
