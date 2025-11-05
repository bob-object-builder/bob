package table

type Property string

const (
	UniqueKey        Property = "unique"
	IndexKey         Property = "index"
	PrimaryKey       Property = "primary"
	OptionalKey      Property = "optional"
	DefaultKey       Property = "="
	AutoIncrementKey Property = "auto_increment"
)

type Column struct {
	name          string
	Type          Type
	Index         bool
	Primary       bool
	Unique        bool
	Optional      bool
	AutoIncrement bool
	Default       string
}

func (c *Column) SetName(name string) {
	c.name = name
}

func (c *Column) GetName() string {
	return c.name
}
