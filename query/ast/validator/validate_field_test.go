package validator

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that the Field aggregation method is valid.
func TestValidateFieldAggregation(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "integer")},
		Statements: ast.Statements{
			&ast.Selection{
				Fields: ast.Fields{
					&ast.Field{Aggregation: "bad_aggr"},
				},
			},
		},
	})
	assert.Equal(t, "field: unsupported aggregation type: bad_aggr", err.Error())
}

// Ensure that the Field DISTINCT is not allowed (for now).
func TestValidateFieldDistinct(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "integer")},
		Statements: ast.Statements{
			&ast.Selection{
				Fields: ast.Fields{
					&ast.Field{Aggregation: "count", Distinct: true},
				},
			},
		},
	})
	assert.Equal(t, "field: DISTINCT is currently unsupported", err.Error())
}

// Ensure that a count() does not have an expression.
func TestValidateFieldCountWithExpression(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "integer")},
		Statements: ast.Statements{
			&ast.Selection{
				Fields: ast.Fields{
					&ast.Field{Aggregation: "count", Expression: &ast.VarRef{Name: "foo"}},
				},
			},
		},
	})
	assert.Equal(t, "field: count() does not accept an expression", err.Error())
}

// Ensure that an aggregation has an appropriate expression.
func TestValidateFieldExpressionType(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "boolean")},
		Statements: ast.Statements{
			&ast.Selection{
				Fields: ast.Fields{
					&ast.Field{Aggregation: "sum", Expression: &ast.VarRef{Name: "foo"}},
				},
			},
		},
	})
	assert.Equal(t, "field: invalid expression data type: boolean", err.Error())
}
