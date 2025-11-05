package table

import (
	"fmt"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/lib/value/object"
	"strings"
)

const Key = "table"

type Table struct {
	name       string
	Columns    object.Object[Column]
	References array.Array[Reference]
}

func New() *Table {
	return &Table{}
}

func (t *Table) SetName(name string) {
	t.name = formatter.ToSnakeCase(name)
}

func (t *Table) GetName() string {
	return t.name
}

func (t *Table) IsNameEmpty() bool {
	return t.name == ""
}

func (t *Table) AddColumn(name string, kind string, properties []string) error {
	column := Column{
		name:          formatter.ToSnakeCase(name),
		Type:          "",
		Default:       "",
		Index:         false,
		Primary:       false,
		Unique:        false,
		Optional:      false,
		AutoIncrement: false,
	}

	hasDefaultValue := false

	typeValue, typeValueError := typeMap.GetType(Type(kind))
	if typeValueError != nil {
		return typeValueError
	}

	column.Type = typeValue

	if typeValue == IdType {
		column.Primary = true
		column.AutoIncrement = true
	}

	for i, token := range properties {
		switch t := Property(token); t {
		case DefaultKey:
			hasDefaultValue = true
			continue
		case UniqueKey:
			column.Unique = true
			continue
		case IndexKey:
			column.Index = true
			continue
		case PrimaryKey:
			column.Primary = true
			continue
		case OptionalKey:
			column.Optional = true
			continue
		default:
			if hasDefaultValue {
				column.Default = formatter.String(strings.Join(properties[i:], " "))
			} else {
				return fmt.Errorf("invalid property '%s'", token)
			}
		}
	}

	t.Columns.Add(name, column)
	return nil
}
