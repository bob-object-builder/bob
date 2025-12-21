package driver

type Literal string

const (
	DateLiteral      Literal = "@date"
	TimeLiteral      Literal = "@time"
	TimestampLiteral Literal = "@now"
)

func GetLiteral(token string) string {
	switch Literal(token) {
	case DateLiteral:
		return "CURRENT_DATE"
	case TimeLiteral:
		return "CURRENT_TIME"
	case TimestampLiteral:
		return "CURRENT_TIMESTAMP"
	}

	return ""
}

func IsLiteral(token string) bool {
	switch Literal(token) {
	case DateLiteral, TimeLiteral, TimestampLiteral:
		return true
	}

	return false
}
