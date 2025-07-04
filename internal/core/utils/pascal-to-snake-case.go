package utils

import (
	"unicode"
)

func PascalToSnakeCase(input string) string {
	var result []rune

	for i, r := range input {

		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
