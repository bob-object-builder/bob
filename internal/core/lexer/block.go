package lexer

import "slices"

type Block struct {
	Command  Command
	actions  []string
	children []any
}

type Blocks []Block

func (b *Block) Append(toAppend any) {
	switch v := toAppend.(type) {
	case Instruction, Block:
		b.children = append(b.children, toAppend)
	case string:
		b.actions = append(b.actions, v)
	}
}

func (b *Block) Actions() []string {
	return b.actions
}

func (b *Block) Children() []any {
	return b.children
}

func (b *Block) ActionIs(toBe ...Command) bool {
	return slices.Contains(toBe, b.Command)
}
