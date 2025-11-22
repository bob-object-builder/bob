package lexer

import (
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/models/order"
)

func (l *Lexer) ParseOrder(target string) (*failure.Failure, *order.Order) {
	l.tokens = l.tokens[1:]
	count := len(l.tokens)

	var nullFirst bool

	if count > 2 {

		if l.tokens[count-2] != order.OrderNullKey {
			return failure.InvalidEmptyNullPriority, nil
		}

		nullFirst = l.tokens[count-1] == order.OrderNullFirst
		l.tokens = l.tokens[:1]
	}

	l.NextLine()

	return nil, &order.Order{
		Direction: l.token,
		Target:    l.ParseReferences(target),
		NullFirst: nullFirst,
	}
}
