package checker

import (
	"salvadorsru/bob/internal/core/kw"
	"strings"
)

func IsExpression(token string) bool {
	return strings.Contains(token, kw.OpenExpression)
}
