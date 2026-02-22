package test

import "testing"

func NormalizeTestString(s string) string {
	result := make([]rune, 0, len(s))
	for _, r := range s {
		if r != ' ' && r != '\n' && r != '\t' {
			result = append(result, r)
		}
	}
	return string(result)
}

func AssertStringsEqual(t *testing.T, driverName, expected, actual string) {
	if NormalizeTestString(actual) != NormalizeTestString(expected) {
		t.Errorf(
			"\nINFO:\ndriver: '%s'\n \nEXPECTED:\n%s\n\nIN:\n%s\n\n",
			driverName,
			expected,
			actual,
		)
	}
}
