package drivers

type Literal string

const (
	Now          = "@now"
	CurrentDate  = "@date"
	CurrentTime  = "@time"
	UtcTimestamp = "@utc"
	SysDate      = "@sysdate"
)
