package mapper

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

// Ensure that a variable declaration can be found in the symbol table.
func TestSymtableFind(t *testing.T) {
	tbl := newSymtable(nil)
	tbl.add(ast.NewVarDecl("foo", "string"))
	tbl.add(ast.NewVarDecl("bar", "integer"))
	if decl := tbl.find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.Name, "foo")
	}
}

// Ensure that a variable declaration can be found in a parent symbol table.
func TestSymtableFindInParentScope(t *testing.T) {
	parent := newSymtable(nil)
	tbl := newSymtable(parent)
	parent.add(ast.NewVarDecl("foo", "string"))
	if decl := tbl.find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.Name, "foo")
	}
}

// Ensure that a variable declaration can be found in a parent symbol table.
func TestSymtableFindOverriddenVariable(t *testing.T) {
	parent := newSymtable(nil)
	tbl := newSymtable(parent)
	parent.add(ast.NewVarDecl("foo", "string"))
	tbl.add(ast.NewVarDecl("foo", "integer"))
	if decl := tbl.find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.DataType, "integer")
	}
}
