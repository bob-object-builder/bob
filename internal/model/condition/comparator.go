package condition

type Comparator string

const (
	And        Comparator = "&&"
	Else       Comparator = "||"
	Like       Comparator = "like"
	Equal      Comparator = "="
	BiggerThan Comparator = ">"
	LowerThan  Comparator = "<"
	Different  Comparator = "!="
)

func IsComparator(token string) bool {
	switch Comparator(token) {
	case And, Else, Like, Equal, BiggerThan, LowerThan, Different:
		return true
	default:
		return false
	}
}
