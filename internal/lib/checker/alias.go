package checker

import "strings"

func IsAlias(name string) bool {
	return strings.HasSuffix(strings.TrimSpace(name), ":")
}
