package ast

import (
	"fmt"
)

// Debug represents a debugging statement in the query.
type Debug struct {
	Expression Expression
}

func (d *Debug) node() {}
func (d *Debug) statement() {}

// NewDebug creates a new debug statement.
func NewDebug() *Debug {
	return &Debug{}
}

func (d *Debug) String() string {
	return fmt.Sprintf("DEBUG(%s)", d.Expression.String())
}
