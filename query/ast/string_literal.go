package ast

import (
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/db"
	"strconv"
)

type StringLiteral struct {
	Value string
}

func (l *StringLiteral) node() string {}

func (l *StringLiteral) String() string {
	return strconv.Quote(l.Value)
}
