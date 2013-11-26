package ast

import "strconv"

// IntegerLiteral represents a hardcoded integer number.
type IntegerLiteral struct {
	Value int
}

// NewIntegerLiteral creates a new IntegerLiteral instance.
func NewIntegerLiteral(value int) *IntegerLiteral {
	return &IntegerLiteral{Value: value}
}

func (l *IntegerLiteral) node()       {}
func (l *IntegerLiteral) expression() {}

func (l *IntegerLiteral) String() string {
	return strconv.Itoa(l.Value)
}
