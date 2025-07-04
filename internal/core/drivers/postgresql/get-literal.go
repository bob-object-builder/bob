package postgresql

import "salvadorsru/bob/internal/core/drivers"

func GetLiteral(toGet string) string {
	switch drivers.Literal(toGet) {
	case drivers.Now:
		return "NOW()"
	case drivers.CurrentDate:
		return "CURRENT_DATE"
	case drivers.CurrentTime:
		return "CURRENT_TIME"
	case drivers.LocalTime:
		// PostgreSQL asume la zona local en CURRENT_TIME
		return "CURRENT_TIME"
	case drivers.LocalTimestamp:
		return "CURRENT_TIMESTAMP"
	case drivers.UtcTimestamp:
		return "CURRENT_TIMESTAMP AT TIME ZONE 'UTC'"
	case drivers.SysDate:
		return "NOW()"
	default:
		return ""
	}
}
