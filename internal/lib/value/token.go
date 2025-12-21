package value

type TokenType int

const (
	Key TokenType = iota
	Expression
	Text
	Alias
)

type Token struct {
	Type  TokenType
	Value string
}
