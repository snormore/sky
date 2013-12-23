package validator

import (
	"fmt"

	"github.com/skydb/sky/query/ast"
)

// Error represents a validation error on an AST node.
type Error struct {
	node ast.Node
	s    string
}

// NewError returns a new error instance.
func NewError(node ast.Node, s string) *Error {
	return &Error{
		node: node,
		s:    s,
	}
}

// errorf returns a validation error using a formatted error string.
func errorf(n ast.Node, format string, args ...interface{}) *Error {
	return NewError(n, fmt.Sprintf(format, args...))
}

// Node returns the AST node that caused the validation error.
func (e *Error) Node() ast.Node {
	return e.node
}

// Error returns the error string.
func (e *Error) Error() string {
	return e.s
}
