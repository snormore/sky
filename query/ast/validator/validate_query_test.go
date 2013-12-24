package validator

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that the Query does not mix INTO with non-INTO selections.
func TestValidateQueryHeterogeneousNaming(t *testing.T) {
	err := Validate(&ast.Query{
		Statements: ast.Statements{
			&ast.Selection{Name:"xxx"},
			&ast.Selection{},
		},
	})
	assert.Equal(t, "query: mixing named and unnamed selections is not allowed", err.Error())
}

// Ensure that the Query does not mix aggregated selections with non-aggregated selections.
func TestValidateQueryHeterogeneousAggregation(t *testing.T) {
	err := Validate(&ast.Query{
		Statements: ast.Statements{
			&ast.Selection{
				Fields: ast.Fields{
					&ast.Field{Aggregation:"count"},
				},
			},
			&ast.Selection{
				Fields: ast.Fields{
					&ast.Field{Expression:&ast.IntegerLiteral{100}},
				},
			},
		},
	})
	assert.Equal(t, "query: mixing aggregated and non-aggregated selection is not allowed", err.Error())
}
