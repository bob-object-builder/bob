package checker

import "strings"

func IsAlias(token string) bool {
	return strings.HasSuffix(token, ":")
}
