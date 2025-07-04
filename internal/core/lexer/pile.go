package lexer

import "fmt"

type Pile map[int]*Block

func NewPile() Pile {
	pile := Pile{}
	pile[0] = &Block{}
	return pile
}

func (p *Pile) Ref(i int) *Block {
	return (*p)[i]
}

func (p *Pile) Add(command string) int {
	i := len(*p)
	(*p)[i] = &Block{Command: Command(command)}
	return i
}

func (p *Pile) MergeLast() bool {
	i := len(*p) - 1
	if i > 1 {
		(*p)[i-1].Append(*(*p)[i])
		delete(*p, i)
		return true
	} else {
		return false
	}
}

func (p *Pile) GetLast() *Block {
	return (*p)[len(*p)-1]
}

func (p *Pile) Clean() {
	for k := range *p {
		delete(*p, k)
	}
}

func (p *Pile) Print() {
	for _, block := range *p {
		fmt.Printf("%s\n", block)
	}
}
