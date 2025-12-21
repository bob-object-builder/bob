package kw

const (
	NewLine         string = "\n"
	Table           string = "table"
	Get             string = "get"
	New             string = "new"
	Delete          string = "delete"
	Raw             string = "raw"
	Drop            string = "drop"
	Set             string = "set"
	Limit           string = "limit"
	Join            string = "->"
	OpenContext     string = "{"
	CloseContext    string = "}"
	VoidContext     string = "{}"
	OpenExpression  string = "("
	CloseExpression string = ")"
	If              string = "if"
	Or              string = "or"
	And             string = "and"
	Group           string = "group"
	Asc             string = "asc"
	Desc            string = "desc"
	Alias           string = ":"
	Nulls           string = "nulls"
	NullsFirst      string = "first"
	NullsLast       string = "last"
	AndSymbol       string = "&&"
	OrSymbol        string = "||"
	// Functions
	ConcatFunction   string = "concat"
	AvgFunction      string = "avg"
	SumFunction      string = "sum"
	MinFunction      string = "min"
	MaxFunction      string = "max"
	CountFunction    string = "count"
	LengthFunction   string = "length"
	OptionalFunction string = "optional"
	Equal            string = "="
)

func IsFunction(token string) bool {
	switch token {
	case ConcatFunction, AvgFunction, SumFunction, MinFunction, MaxFunction, CountFunction, LengthFunction:
		return true
	default:
		return false
	}
}

func IsKeyword(token string) bool {
	switch token {
	case NewLine, Table, Get, New, Delete, Raw, Drop, Set, Join, OpenContext, CloseContext, VoidContext, OpenExpression, CloseExpression, If, Or, And, Group, Asc, Desc, Alias, Nulls, NullsFirst, NullsLast, AndSymbol, OrSymbol:
		return true
	default:
		return false
	}
}
