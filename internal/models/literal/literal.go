package literal

import "strings"

const DateKey = "@date"
const TimeKey = "@time"
const TimestampKey = "@now"
const NullKey = "null"

func IsLiteral(token string) bool {
	return strings.HasPrefix(token, "@")
}

func GetLiteral(token string) string {
	switch token {
	case DateKey:
		return "CURRENT_DATE"
	case TimeKey:
		return "CURRENT_TIME"
	case TimestampKey:
		return "CURRENT_TIMESTAMP"
	}

	return token
}
