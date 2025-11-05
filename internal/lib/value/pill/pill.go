package pill

type Pill string

func (a *Pill) Set(token string) {
	*a = Pill(token)
}

func (a *Pill) Use() string {
	token := string(*a)
	*a = ""
	return token
}

func (a *Pill) UseOr(token string) string {
	if a.IsEmpty() {
		return token
	}
	return a.Use()
}

func (p *Pill) IsEmpty() bool {
	return string(*p) == ""
}

func New(alias string) Pill {
	return Pill(alias)
}
