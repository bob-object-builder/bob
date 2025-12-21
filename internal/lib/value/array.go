package value

import "fmt"

type Array[T any] []T

func (l *Array[T]) Push(toAppend ...T) {
	*l = append(*l, toAppend...)
}

func (l *Array[T]) Prepend(toPrepend ...T) {
	*l = append(toPrepend, *l...)
}

func (l *Array[T]) GetLast() *T {
	if len(*l) == 0 {
		return nil
	}
	return &(*l)[len(*l)-1]
}

func (l *Array[T]) Get(i int) *T {
	if i < 0 || i >= len(*l) {
		return nil
	}
	return &(*l)[i]
}

func NewArray[T any](values ...T) *Array[T] {
	if len(values) > 0 {
		arr := Array[T](values)
		return &arr
	}

	return &Array[T]{}
}

func (l *Array[T]) Pop() *T {
	if len(*l) == 0 {
		return nil
	}

	last := (*l)[len(*l)-1]
	*l = (*l)[:len(*l)-1]
	return &last
}

func (l *Array[T]) Clean() {
	*l = Array[T](nil)
}

func (l *Array[T]) Length() int {
	return len(*l)
}

func (l *Array[T]) Slice(quantity int) *Array[T] {
	if quantity <= 0 {
		return NewArray[T]()
	}

	if len(*l) < quantity {
		quantity = len(*l)
	}

	newArr := Array[T]((*l)[:quantity])
	return &newArr
}

func (l *Array[T]) Join(separator string) string {
	if len(*l) == 0 {
		return ""
	}

	var result string
	for i, value := range *l {
		result += fmt.Sprintf("%v", value)
		if i < len(*l)-1 {
			result += separator
		}
	}
	return result
}
