package lexer

type Command string

const (
	Table         Command = "table"
	Get           Command = "get"
	LeftJoin      Command = "left"
	LeftJoinAlias Command = "->"
	New           Command = "new"
	Set           Command = "set"
	Delete        Command = "delete"
)

type Direction string

const (
	Left Direction = Direction(LeftJoin)
)

func IsCommand(s string) bool {
	switch Command(s) {
	case Table, Get, LeftJoin, LeftJoinAlias, New, Set, Delete:
		return true
	default:
		return false
	}
}
