package ast

import (
	"strconv"
)

type StringLiteral struct {
	Value string
}

// NewStringLiteral creates a new StringLiteral instance.
func NewStringLiteral(value string) *StringLiteral {
	return &StringLiteral{Value: value}
}

func (l *StringLiteral) node()       {}
func (l *StringLiteral) expression() {}

func (l *StringLiteral) String() string {
	return strconv.Quote(l.Value)
}
