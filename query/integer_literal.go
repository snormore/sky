package query

import "strconv"

type IntegerLiteral struct {
	queryElementImpl
	value int
}

func (l *IntegerLiteral) String() string {
	return strconv.Itoa(l.value)
}
