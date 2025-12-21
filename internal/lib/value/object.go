package value

import (
	"fmt"
	"strings"
)

type child[T any] struct {
	index int
	Value T
}

type Map[T any] struct {
	length int
	Values map[string]child[T]
	order  []string
}

func NewMap[T any]() *Map[T] {
	return &Map[T]{
		length: 0,
		Values: make(map[string]child[T]),
		order:  make([]string, 0),
	}
}

func (obj *Map[T]) Get(id string) *T {
	if c, ok := obj.Values[id]; ok {
		return &c.Value
	}
	return nil
}

func (obj *Map[T]) Set(id string, value T) {
	if obj.Values == nil {
		obj.Values = make(map[string]child[T])
		obj.order = make([]string, 0)
	}

	if _, exists := obj.Values[id]; !exists {
		obj.order = append(obj.order, id)
	}

	obj.Values[id] = child[T]{
		index: obj.length,
		Value: value,
	}
	obj.length++
}

func (obj *Map[T]) Range() <-chan struct {
	Key   string
	Value T
} {
	ch := make(chan struct {
		Key   string
		Value T
	})
	go func() {
		for _, key := range obj.order {
			c := obj.Values[key]
			ch <- struct {
				Key   string
				Value T
			}{key, c.Value}
		}
		close(ch)
	}()
	return ch
}

func (obj *Map[T]) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for _, key := range obj.order {
		if !first {
			sb.WriteString(", ")
		}
		c := obj.Values[key]
		sb.WriteString(fmt.Sprintf("%q: %v", key, c.Value))
		first = false
	}
	sb.WriteString("}")
	return sb.String()
}
