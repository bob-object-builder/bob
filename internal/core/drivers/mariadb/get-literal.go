package mariadb

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
		return "CURRENT_TIME"
	case drivers.LocalTimestamp:
		return "CURRENT_TIMESTAMP"
	case drivers.UtcTimestamp:
		return "UTC_TIMESTAMP()"
	case drivers.SysDate:
		return "SYSDATE()"
	default:
		return ""
	}
}
