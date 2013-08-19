package query

import "strconv"

type IntegerLiteral struct {
	queryElementImpl
	value int
}

func (l *IntegerLiteral) Codegen() (string, error) {
	return l.String(), nil
}

func (l *IntegerLiteral) String() string {
	return strconv.Itoa(l.value)
}
