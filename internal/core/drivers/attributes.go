package drivers

type Attributes struct {
	Primary       string
	Required      string
	Index         string
	AutoIncrement string
	Unique        string
	Nullable      string
}

type Attribute string

func (t Attributes) Get(attr Attribute) (string, bool) {
	switch attr {
	case Primary:
		return t.Primary, true
	case Required:
		return t.Required, true
	case Index:
		return t.Index, true
	case AutoIncrement:
		return t.AutoIncrement, true
	case Unique:
		return t.Unique, true
	case Nullable:
		return t.Nullable, true
	default:
		return "", false
	}
}

const (
	Primary       Attribute = "primary"
	Required      Attribute = "required"
	Index         Attribute = "index"
	AutoIncrement Attribute = "auto_increment"
	Unique        Attribute = "unique"
	Nullable      Attribute = "nullable"
	Default       Attribute = "=" // Default value
	ForeignKey    Attribute = "foreign_key"
	Check         Attribute = "check"     // CHECK constraints
	OnUpdate      Attribute = "on_update" // ON UPDATE actions
	OnDelete      Attribute = "on_delete" // ON DELETE actions
	Unsigned      Attribute = "unsigned"  // Only for unsigned integers
	Generated     Attribute = "generated" // Generated columns (MySQL)
	Computed      Attribute = "computed"  // Computed columns
	Collation     Attribute = "collation" // For text
	Comment       Attribute = "comment"   // Associated comment
)
