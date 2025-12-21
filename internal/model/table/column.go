package table

import (
	"salvadorsru/bob/internal/core/drivers/driver"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/kw"
)

type Column struct {
	Name            string
	Type            string
	Default         string
	NextIsDefault   bool
	IsOptional      bool
	IsReference     bool
	IsAutoIncrement bool
	IsCurrent       bool
	IsPrimary       bool
	IsIndex         bool
	IsUnique        bool
}

func NewColumn(name string) *Column {
	return &Column{
		Name: name,
	}
}

func (c *Column) Parse(token string) *failure.Failure {
	if token == kw.NewLine {
		return nil
	}

	if c.Name == "" {
		c.Name = token
		return nil
	}

	if c.Type == "" {
		if driver.IsType(token) {
			c.Type = token
			return nil
		} else {
			return failure.UndefinedTypeForColumn(token, c.Name)
		}
	}

	switch Property(token) {
	case Index:
		c.IsIndex = true
		return nil
	case Unique:
		c.IsUnique = true
		return nil
	case Primary:
		c.IsPrimary = true
		return nil
	case Optional:
		c.IsOptional = true
		return nil
	}

	if token == kw.Equal {
		c.NextIsDefault = true
		return nil
	}

	if c.NextIsDefault {
		if driver.IsLiteral(token) {
			c.Default = driver.GetLiteral(token)
			return nil
		}

		c.Default = token
		return nil
	}

	return nil
}
