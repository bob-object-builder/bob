package sqlite

import "salvadorsru/bob/internal/core/drivers"

func GetLiteral(toGet string) string {
	switch drivers.Literal(toGet) {
	case drivers.Now:
		return "datetime('now')"
	case drivers.CurrentDate:
		return "date('now')"
	case drivers.CurrentTime:
		return "time('now')"
	case drivers.LocalTime:
		return "time('now', 'localtime')"
	case drivers.LocalTimestamp:
		return "datetime('now', 'localtime')"
	case drivers.UtcTimestamp:
		return "datetime('now')"
	case drivers.SysDate:
		return "datetime('now')"
	default:
		return ""
	}
}
