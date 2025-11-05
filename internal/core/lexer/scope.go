package lexer

import "strings"

type ScopeKey string

const (
	OpenKey  ScopeKey = "{"
	CloseKey ScopeKey = "}"
)

func IsOpenKey(token string) bool {
	return strings.HasPrefix(token, string(OpenKey))
}

func IsCloseKey(token string) bool {
	return strings.HasSuffix(token, string(CloseKey))
}
