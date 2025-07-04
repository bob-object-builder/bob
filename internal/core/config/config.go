package config

type Model struct {
	ActionAlias string
}

var config *Model = &Model{
	ActionAlias: ".",
}

func New(newConfig Model) {
	config = &newConfig
}

func Use() Model {
	if config == nil {
		panic("undefined: config")
	}

	return *config
}
