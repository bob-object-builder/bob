package table

import (
	"fmt"
	"salvadorsru/bob/internal/lib/formatter"
	"strings"
)

const (
	OnDeleteCascadeKey Property = "cascade"
	OnUpdateCascadeKey Property = "propagate"
	BindKey            Property = "bind"
)

type Reference struct {
	Table           string
	Column          string
	Optional        bool
	OnDeleteCascade bool
	OnUpdateCascade bool
	Default         string
}

func (t *Table) AddReference(table string, column string, properties []string) error {
	ref := Reference{
		Table:           formatter.ToSnakeCase(table),
		Column:          formatter.ToSnakeCase(column),
		Optional:        false,
		OnDeleteCascade: false,
		OnUpdateCascade: false,
	}

	hasDefaultValue := false

	for i, token := range properties {
		switch t := Property(token); t {
		case DefaultKey:
			hasDefaultValue = true
			continue
		case OptionalKey:
			ref.Optional = true
			continue
		case OnDeleteCascadeKey:
			ref.OnDeleteCascade = true
			continue
		case OnUpdateCascadeKey:
			ref.OnUpdateCascade = true
		case BindKey:
			ref.OnDeleteCascade = true
			ref.OnUpdateCascade = true

		default:
			if hasDefaultValue {
				ref.Default = formatter.NormalizeString(strings.Join(properties[i:], " "))
			} else {
				return fmt.Errorf("invalid property '%s'", token)
			}
		}
	}

	t.References.Push(ref)
	return nil
}
