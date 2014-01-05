package validator

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that a BinaryExpression has matching types on both sides.
func TestValidateBinaryExpressionDataTypeMismatch(t *testing.T) {
	err := Validate(&ast.BinaryExpression{LHS: &ast.IntegerLiteral{10}, RHS: &ast.BooleanLiteral{true}, Op: ast.OpEquals})
	assert.Equal(t, "expression: data type mismatch: integer != boolean", err.Error())
}

// Ensure that a boolean BinaryExpression has an appropriate operator.
func TestValidateBinaryExpressionInvalidBooleanOperator(t *testing.T) {
	err := Validate(&ast.BinaryExpression{LHS: &ast.BooleanLiteral{true}, RHS: &ast.BooleanLiteral{true}, Op: ast.OpPlus})
	assert.Equal(t, "expression: invalid boolean operator: +", err.Error())
}

// Ensure that a factor BinaryExpression has an appropriate operator.
func TestValidateBinaryExpressionInvalidFactorOperator(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "factor")},
		Statements: ast.Statements{
			&ast.Condition{
				Expression: &ast.BinaryExpression{
					LHS: &ast.VarRef{Name: "foo"},
					RHS: &ast.StringLiteral{"XXX"},
					Op:  ast.OpPlus,
				},
			},
		},
	})
	assert.Equal(t, "expression: invalid factor operator: +", err.Error())
}

// Ensure that a factor BinaryExpression has matching LHS and RHS factors.
func TestValidateBinaryExpressionFactorAssociationMismatch(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{
			ast.NewVarDecl(0, "foo", "factor"),
			&ast.VarDecl{Name: "bar", DataType: "factor"},
		},
		Statements: ast.Statements{
			&ast.Condition{
				Expression: &ast.BinaryExpression{
					LHS: &ast.VarRef{Name: "foo"},
					RHS: &ast.VarRef{Name: "bar"},
					Op:  ast.OpPlus,
				},
			},
		},
	})
	assert.Equal(t, "expression: mismatched factor association: foo<> != bar<>", err.Error())
}

// Ensure that a factor BinaryExpression does not have two strings.
func TestValidateBinaryExpressionDoubleStringLiterals(t *testing.T) {
	err := Validate(&ast.Condition{
		Expression: &ast.BinaryExpression{
			LHS: &ast.StringLiteral{"foo"},
			RHS: &ast.StringLiteral{"bar"},
			Op:  ast.OpEquals,
		},
	})
	assert.Equal(t, "expression: string literal comparison not allowed: \"foo\" == \"bar\"", err.Error())
}

// Ensure that a factor BinaryExpression with matching associations validates.
func TestValidateBinaryExpressionFactorAssociation(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{
			ast.NewVarDecl(0, "foo", "factor"),
			&ast.VarDecl{Name: "bar", DataType: "factor", Association: "foo"},
		},
		Statements: ast.Statements{
			&ast.Condition{
				Expression: &ast.BinaryExpression{
					LHS: &ast.VarRef{Name: "foo"},
					RHS: &ast.VarRef{Name: "bar"},
					Op:  ast.OpNotEquals,
				},
			},
		},
	})
	assert.NoError(t, err)
}

// Ensure that a integer BinaryExpression has an appropriate operator.
func TestValidateBinaryExpressionInvalidIntegerOperator(t *testing.T) {
	err := Validate(&ast.BinaryExpression{LHS: &ast.IntegerLiteral{100}, RHS: &ast.IntegerLiteral{100}, Op: ast.OpAnd})
	assert.Equal(t, "expression: invalid numeric operator: &&", err.Error())
}

// Ensure that a BinaryExpression is valid.
func TestValidateBinaryExpression(t *testing.T) {
	err := Validate(&ast.BinaryExpression{LHS: &ast.IntegerLiteral{10}, RHS: &ast.IntegerLiteral{10}, Op: ast.OpEquals})
	assert.NoError(t, err)
}
