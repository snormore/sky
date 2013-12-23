package validator

import (
	"github.com/skydb/sky/query/ast"
)

// Validator is used to ensure that an AST tree can be used
// by the mapper and reducer.
type Validator interface {
	Validate(ast.Node) error
}

// validator is the concrete implementation of the Validator interface.
type validator struct {
	stack     []ast.Node
	dataTypes map[ast.Node]string
	err       error
}

// New creates a new Validator instance.
func New() Validator {
	return &validator{}
}

// Validate validates a given AST node.
func Validate(n ast.Node) error {
	return New().Validate(n)
}

// Validate walks the AST and validates every node and returns the errors.
func (v *validator) Validate(n ast.Node) error {
	v.stack = make([]ast.Node, 0)
	v.dataTypes = make(map[ast.Node]string)
	v.err = nil
	ast.Walk(v, n)
	return v.err
}

// Visit implements the Visitor interface.
func (v *validator) Visit(n ast.Node, tbl *ast.Symtable) ast.Visitor {
	v.stack = append(v.stack, n)

	switch n := n.(type) {
	case *ast.Assignment:
		v.visitAssignment(n, tbl)
	case *ast.BooleanLiteral:
		v.visitBooleanLiteral(n, tbl)
	case *ast.IntegerLiteral:
		v.visitIntegerLiteral(n, tbl)
	case *ast.StringLiteral:
		v.visitStringLiteral(n, tbl)
	case *ast.VarRef:
		v.visitVarRef(n, tbl)
	}

	// Stop walking once we hit an error.
	if v.err != nil {
		return nil
	}
	return v
}

// Visit implements the Visitor interface.
func (v *validator) Exiting(n ast.Node, tbl *ast.Symtable) {
	v.stack = v.stack[0 : len(v.stack)-1]

	if v.err != nil {
		return
	}

	switch n := n.(type) {
	case *ast.Assignment:
		v.exitingAssignment(n, tbl)
	case *ast.BinaryExpression:
		v.exitingBinaryExpression(n, tbl)
	}
}

// Visit implements the Visitor interface.
func (v *validator) Error(err error) ast.Visitor {
	v.err = NewError(nil, err.Error())
	return nil
}
