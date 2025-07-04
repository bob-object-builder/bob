package lexer

type Tables struct {
	List  map[string]Block
	Order []string
}

func NewTables() Tables {
	return Tables{
		List:  make(map[string]Block),
		Order: []string{},
	}
}

func (t *Tables) Set(key string, value Block) {
	if _, exists := t.List[key]; !exists {
		t.Order = append(t.Order, key)
	}
	t.List[key] = value
}

func (t *Tables) Get(key string) (Block, bool) {
	val, ok := t.List[key]
	return val, ok
}

func (t *Tables) Keys() []string {
	return t.Order
}

func (t *Tables) Values() []Block {
	values := make([]Block, 0, len(t.Order))
	for _, key := range t.Order {
		values = append(values, t.List[key])
	}
	return values
}

type Program struct {
	Tables  Tables
	Actions Blocks
}

func NewProgram() *Program {
	return &Program{
		Tables:  NewTables(),
		Actions: make(Blocks, 0),
	}
}

func (p *Program) Append(block Block) {
	if block.Command == Table {
		tableRef := block.Actions()[0]
		p.Tables.Set(tableRef, block)
		return
	}

	p.Actions = append(p.Actions, block)
}
