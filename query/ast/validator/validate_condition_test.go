package validator

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that the Condition UOM is valid.
func TestValidateAssignmentUOM(t *testing.T) {
	err := Validate(&ast.Condition{UOM: "sessions"})
	assert.Equal(t, "condition: unit not supported: sessions", err.Error())
}

// Ensure that the Condition start is not before the end.
func TestValidateAssignmentStartEnd(t *testing.T) {
	err := Validate(&ast.Condition{Start: 3, End: 2})
	assert.Equal(t, "condition: start index (3) cannot be less than end index (2)", err.Error())
}

// Ensure that the Condition return type is boolean
func TestValidateAssignmentExpressionDataType(t *testing.T) {
	err := Validate(&ast.Condition{
		Expression: &ast.BinaryExpression{
			LHS: &ast.IntegerLiteral{100},
			RHS: &ast.IntegerLiteral{200},
			Op:  ast.OpPlus,
		},
	})
	assert.Equal(t, "condition: invalid expression return type: integer", err.Error())
}
