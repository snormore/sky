package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure that a variable declaration can be found in the symbol table.
func TestSymtableFind(t *testing.T) {
	tbl := NewSymtable(nil)
	tbl.Add(NewVarDecl(0, "foo", "string"))
	tbl.Add(NewVarDecl(0, "bar", "integer"))
	if decl := tbl.Find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.Name, "foo")
	}
}

// Ensure that a variable declaration can be found in a parent symbol table.
func TestSymtableFindInParentScope(t *testing.T) {
	parent := NewSymtable(nil)
	tbl := NewSymtable(parent)
	parent.Add(NewVarDecl(0, "foo", "string"))
	if decl := tbl.Find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.Name, "foo")
	}
}

// Ensure that a variable declaration can be found in a parent symbol table.
func TestSymtableFindOverriddenVariable(t *testing.T) {
	parent := NewSymtable(nil)
	tbl := NewSymtable(parent)
	parent.Add(NewVarDecl(0, "foo", "string"))
	tbl.Add(NewVarDecl(0, "foo", "integer"))
	if decl := tbl.Find("foo"); assert.NotNil(t, decl) {
		assert.Equal(t, decl.DataType, "integer")
	}
}
