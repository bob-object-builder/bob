package table

import "salvadorsru/bob/internal/lib/formatter"

type Reference struct {
	Table      string
	Column     string
	IsOptional bool
}

func (t *Table) AddReference(table string, column string, isOptional bool) {
	t.References.Push(
		Reference{
			Table:      formatter.ToSnakeCase(table),
			Column:     formatter.ToSnakeCase(column),
			IsOptional: isOptional,
		},
	)
}
