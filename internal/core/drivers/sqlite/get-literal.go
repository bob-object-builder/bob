package sqlite

import "salvadorsru/bob/internal/core/drivers"

func GetLiteral(toGet string) string {
	switch drivers.Literal(toGet) {
	case drivers.Now, drivers.UtcTimestamp, drivers.SysDate:
		return "CURRENT_TIMESTAMP"
	case drivers.CurrentDate:
		return "CURRENT_DATE"
	case drivers.CurrentTime:
		return "CURRENT_TIME"
	default:
		return ""
	}
}
