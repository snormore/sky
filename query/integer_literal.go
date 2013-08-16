package query

import "strconv"

type IntegerLiteral struct {
	value int
}

func (l *IntegerLiteral) String() string {
	return strconv.Itoa(l.value)
}
