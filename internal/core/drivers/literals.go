package drivers

type Literal string

const (
	Now            = "@now"       // Current date and time (equivalent to CURRENT_TIMESTAMP, NOW(), datetime('now'), etc)
	CurrentDate    = "@date"      // Current date (CURRENT_DATE, CURDATE(), date('now'))
	CurrentTime    = "@time"      // Current time (CURRENT_TIME, CURTIME(), time('now'))
	LocalTime      = "@localTime" // Local time without date (LOCALTIME, time('now','localtime'))
	LocalTimestamp = "@locaNow"   // Local date and time (LOCALTIMESTAMP, datetime('now','localtime'))
	UtcTimestamp   = "@utc"       // UTC date and time (UTC_TIMESTAMP(), datetime('now'))
	SysDate        = "@sysdate"   // Evaluated at the exact moment (SYSDATE(), CLOCK_TIMESTAMP())
)
