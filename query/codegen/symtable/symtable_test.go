package symtable

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that a variable declaration can be found in the symbol table.
func TestSymtableFind(t *testing.T) {
	tbl := New(nil)
	tbl.Add(ast.NewVarDecl("foo", "string"))
	tbl.Add(ast.NewVarDecl("bar", "integer"))
	if decl := tbl.Find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.Name, "foo")
	}
}

// Ensure that a variable declaration can be found in a parent symbol table.
func TestSymtableFindInParentScope(t *testing.T) {
	parent := New(nil)
	tbl := New(parent)
	parent.Add(ast.NewVarDecl("foo", "string"))
	if decl := tbl.Find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.Name, "foo")
	}
}

// Ensure that a variable declaration can be found in a parent symbol table.
func TestSymtableFindOverriddenVariable(t *testing.T) {
	parent := New(nil)
	tbl := New(parent)
	parent.Add(ast.NewVarDecl("foo", "string"))
	tbl.Add(ast.NewVarDecl("foo", "integer"))
	if decl := tbl.Find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.DataType, "integer")
	}
}
