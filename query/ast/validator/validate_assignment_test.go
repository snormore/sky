package validator

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that an Assignment has a target.
func TestValidateAssignmentTargetRequired(t *testing.T) {
	err := Validate(&ast.Assignment{})
	assert.Equal(t, "assignment: target required", err.Error())
}

// Ensure that an Assignment assigns using the right data type.
func TestValidateAssignmentDataTypeMismatch(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "boolean")},
		Statements:     ast.Statements{&ast.Assignment{Target: &ast.VarRef{"foo"}, Expression: &ast.IntegerLiteral{10}}},
	})
	assert.Equal(t, "assignment: data type mismatch: boolean != integer", err.Error())
}

// Ensure that an Assignment is valid.
func TestValidateAssignment(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "integer")},
		Statements:     ast.Statements{&ast.Assignment{Target: &ast.VarRef{"foo"}, Expression: &ast.IntegerLiteral{10}}},
	})
	assert.NoError(t, err)
}
