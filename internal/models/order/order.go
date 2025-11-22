package order

const OrderAscKey = "asc"
const OrderDescKey = "desc"
const OrderNullKey = "null"
const OrderNullFirst = "first"

type Order struct {
	Direction string
	Target    string
	NullFirst bool
}

func IsOrder(key string) bool {
	return key == OrderAscKey || key == OrderDescKey
}
