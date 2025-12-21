package stack

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/table"
)

type Stack[T any] struct {
	items   value.Array[T]
	tables  value.Map[*table.Table]
	actions value.Array[T]
	output  value.Array[T]
}

type Item[T any] struct {
	Value T
}

func NewStack[T any]() *Stack[T] {
	return &Stack[T]{
		tables:  *value.NewMap[*table.Table](),
		items:   *value.NewArray[T](),
		actions: *value.NewArray[T](),
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

type HasMerge interface {
	Merge(i any) *failure.Failure
}

func (s *Stack[T]) Merge() *failure.Failure {
	if s.items.Length() < 2 {
		popped := s.Pop()
		if popped == nil {
			return nil
		}

		switch i := any(*popped).(type) {
		case *table.Table:
			s.tables.Set(i.Name, i)
		default:
			s.actions.Push(*popped)
		}

		// s.output.Push(*popped)
		return nil
	}

	last := s.items.Pop()
	previous := s.items.GetLast()

	if last == nil || previous == nil {
		return failure.TypeDoesNotImplementMerge
	}

	if merger, ok := any(*previous).(HasMerge); ok {
		return merger.Merge(*last)
	}

	return failure.TypeDoesNotImplementMerge
}

func (s *Stack[T]) GetTables() value.Map[*table.Table] {
	return s.tables
}

func (s *Stack[T]) GetActions() value.Array[T] {
	return s.actions
}

func (s *Stack[T]) GetOutput() value.Array[T] {
	return s.output
}
