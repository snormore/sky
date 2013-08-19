package query

import "strconv"

type StringLiteral struct {
	queryElementImpl
	value string
}

func (l *StringLiteral) String() string {
	return strconv.Quote(l.value)
}
