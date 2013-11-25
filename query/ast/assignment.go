package ast

import (
	"bytes"
	"fmt"
	"github.com/skydb/sky/core"
)

// An Assignment statement sets a variable to a given value.
type Assignment struct {
	Target     *VarRef
	Expression Expression
}

func (a *Assignment) node() string {}

// Creates a new assignment.
func NewAssignment() *Assignment {
	return &Assignment{}
}

// Converts the statement to a string-based representation.
func (a *Assignment) String() string {
	return fmt.Sprintf("SET @%s = %s", a.target.value, a.expression.String())
}
