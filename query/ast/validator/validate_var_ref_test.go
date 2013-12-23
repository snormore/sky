package validator

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that a VarRef exists.
func TestValidateVarRefNotFound(t *testing.T) {
	err := Validate(&ast.VarRef{"foo"})
	assert.Equal(t, "variable: not found: foo", err.Error())
}

// Ensure that a VarRef is valid.
func TestValidateVarRef(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "integer")},
		Statements: ast.Statements{
			&ast.Assignment{Target: &ast.VarRef{"foo"}, Expression: &ast.IntegerLiteral{100}},
		},
	})
	assert.NoError(t, err)
}
