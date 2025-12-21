package lexer

type TokenType int

const (
	Key TokenType = iota
	Word
	Expression
	Text
)

type Token struct {
	Type  TokenType
	Value string
}
