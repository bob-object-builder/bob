package table

import (
	"salvadorsru/bob/internal/core/lexer"
	"salvadorsru/bob/internal/core/utils"
)

type Column struct {
	Type            string
	Attributes      lexer.Instruction
	IsPrimaryKey    bool
	IsAutoIncrement bool
	Default         *string
}

type Table struct {
	Id          string
	Name        string
	Columns     utils.Object[*Column]
	Indexes     []Indexable
	Uniques     []string
	PrimaryKeys []string
	References  utils.Object[Reference]
}
