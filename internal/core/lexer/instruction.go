package lexer

type Instruction []string

// NewInstructions creates a new Instructions from a slice of strings.
func NewInstructions(instructions ...string) Instruction {
	return instructions
}

// AppendStep appends a step to the instruction.
func (i *Instruction) Append(instruction string) {
	*i = append(*i, instruction)
}
