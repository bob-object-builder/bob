package utils

type Object[T comparable] struct {
	Data  map[string]T
	Order []string
}

func NewObject[T comparable]() Object[T] {
	return Object[T]{
		Data:  make(map[string]T),
		Order: []string{},
	}
}

func (o *Object[T]) Set(key string, value T) {
	if o.Data == nil {
		o.Data = make(map[string]T)
	}
	if o.Order == nil {
		o.Order = []string{}
	}
	if _, exists := o.Data[key]; !exists {
		o.Order = append(o.Order, key)
	}
	o.Data[key] = value
}

func (o *Object[T]) Get(key string) T {
	val, _ := o.Data[key]
	return val
}

func (o *Object[T]) Merge(other Object[T]) {
	for _, key := range other.Order {
		o.Set(key, other.Data[key])
	}
}

func (o *Object[T]) Keys() []string {
	return o.Order
}

func (o *Object[T]) Values() []T {
	values := make([]T, 0, len(o.Order))
	for _, k := range o.Order {
		values = append(values, o.Data[k])
	}
	return values
}

func (o *Object[T]) Has(key string) bool {
	_, exists := o.Data[key]
	return exists
}

func (o *Object[T]) Reverse() Object[T] {
	newObj := NewObject[T]()
	for i := len(o.Order) - 1; i >= 0; i-- {
		key := o.Order[i]
		newObj.Set(key, o.Data[key])
	}
	return newObj
}

func (o *Object[T]) Prepend(key string, value T) {
	if o.Data == nil {
		o.Data = make(map[string]T)
	}
	if o.Order == nil {
		o.Order = []string{}
	}
	if _, exists := o.Data[key]; !exists {
		o.Order = append([]string{key}, o.Order...)
	}
	o.Data[key] = value
}

func (o *Object[T]) Delete(key string) {
	if _, exists := o.Data[key]; exists {
		delete(o.Data, key)
		newOrder := make([]string, 0, len(o.Order)-1)
		for _, k := range o.Order {
			if k != key {
				newOrder = append(newOrder, k)
			}
		}
		o.Order = newOrder
	}
}
