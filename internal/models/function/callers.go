package function

import (
	"salvadorsru/bob/internal/core/failure"
)

type Caller string

type CallersMap map[Caller]Caller

type Callers struct {
	Length Caller
}

func NewCallers(t Callers) CallersMap {
	return CallersMap{
		"length": t.Length,
	}
}

var callerMap = CallersMap{
	"length": "length",
}

func (t CallersMap) GetCaller(token string) (*failure.Failure, Caller) {
	if typ, ok := t[Caller(token)]; ok {
		return nil, typ
	}
	return failure.UndefinedCaller, ""
}

func IsType(token string) bool {
	if typ, ok := callerMap[Caller(token)]; ok {
		return typ != ""
	}
	return false
}
