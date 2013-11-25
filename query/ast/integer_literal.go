package ast

import "strconv"

// IntegerLiteral represents a hardcoded integer number.
type IntegerLiteral struct {
	Value int
}

func (l *IntegerLiteral) node() string {}

func (l *IntegerLiteral) String() string {
	return strconv.Itoa(l.Value)
}
