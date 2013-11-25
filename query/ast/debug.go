package ast

import (
	"bytes"
	"fmt"
)

// Debug represents a debugging statement in the query.
type Debug struct {
	Expression Expression
}

func (d *Debug) node() string {}

// NewDebug creates a new debug statement.
func NewDebug() *Debug {
	return &Debug{}
}

func (d *Debug) String() string {
	return fmt.Sprintf("DEBUG(%s)", d.expression.String())
}
