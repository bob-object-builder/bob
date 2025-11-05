package lexer

import "salvadorsru/bob/internal/lib/value/array"

type Stack[T any] struct {
	items array.Array[T]
}

type Item[T any] struct {
	Value T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		items: *array.New[T](),
	}
}

func (s *Stack[T]) Push(value T) {
	s.items.Push(value)
}

func (s *Stack[T]) GetLast() *T {
	return s.items.GetLast()
}

func (s *Stack[T]) Length() int {
	return s.items.Length()
}

func (s *Stack[T]) Clean() {
	s.items.Clean()
}

func (s *Stack[T]) Pop() *T {
	return s.items.Pop()
}
