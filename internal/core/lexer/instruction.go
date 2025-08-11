package lexer

import "slices"

type Instruction []string

// NewInstructions creates a new Instructions from a slice of strings.
func NewInstructions(instructions ...string) Instruction {
	return instructions
}

// AppendStep appends a step to the instruction.
func (i *Instruction) Append(instruction string) {
	*i = append(*i, instruction)
}

// Has checks if the instruction contains the given step.
func (i *Instruction) Has(instruction string) bool {
	return slices.Contains(*i, instruction)
}
