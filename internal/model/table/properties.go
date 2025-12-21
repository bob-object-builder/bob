package table

type Property string

const (
	Index         Property = "index"
	Unique        Property = "unique"
	Primary       Property = "primary"
	Optional      Property = "optional"
	AutoIncrement Property = "auto_increment"
)

func IsProperty(s string) bool {
	switch Property(s) {
	case Index, Unique, Primary, Optional, AutoIncrement:
		return true
	default:
		return false
	}
}
