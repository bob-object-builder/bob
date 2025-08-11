package drivers

type Attributes struct {
	Primary       string
	Required      string
	Index         string
	AutoIncrement string
	Unique        string
	Optional      string
}

type Attribute string

func (t Attributes) Get(attr Attribute) (string, bool) {
	switch attr {
	case Primary:
		return t.Primary, true
	case Index:
		return t.Index, true
	case AutoIncrement:
		return t.AutoIncrement, true
	case Unique:
		return t.Unique, true
	case Optional:
		return t.Optional, true
	default:
		return "", false
	}
}

const (
	Primary       Attribute = "primary"
	Index         Attribute = "index"
	AutoIncrement Attribute = "auto_increment"
	Unique        Attribute = "unique"
	Optional      Attribute = "optional"
	Default       Attribute = "=" // Default value
)
