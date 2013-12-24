package validator

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that the Selection dimension is found.
func TestValidateSelectionDimensionExists(t *testing.T) {
	err := Validate(&ast.Selection{Dimensions: []string{"foo"}})
	assert.Equal(t, "selection: dimension variable not found: foo", err.Error())
}

// Ensure that the Selection string dimension is not allowed.
func TestValidateSelectionStringDimension(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "string")},
		Statements: ast.Statements{
			&ast.Selection{Dimensions: []string{"foo"}},
		},
	})
	assert.Equal(t, "selection: string variables cannot be used as dimensions: foo", err.Error())
}

// Ensure that the Selection float dimension is not allowed.
func TestValidateSelectionFloatDimension(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "float")},
		Statements: ast.Statements{
			&ast.Selection{Dimensions: []string{"foo"}},
		},
	})
	assert.Equal(t, "selection: float variables cannot be used as dimensions: foo", err.Error())
}

// Ensure that the Selection fields are not used twice.
func TestValidateSelectionDuplicateFieldIdentifiers(t *testing.T) {
	err := Validate(&ast.Query{
		SystemVarDecls: ast.VarDecls{ast.NewVarDecl(0, "foo", "integer")},
		Statements: ast.Statements{
			&ast.Selection{
				Fields: ast.Fields{
					&ast.Field{Name:"count", Aggregation:"sum", Expression:&ast.VarRef{Name:"foo"}},
					&ast.Field{Aggregation:"count"},
				},
			},
		},
	})
	assert.Equal(t, "selection: field name already used: count", err.Error())
}
