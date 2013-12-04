package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure that basic variable declarations are found.
func TestFindVarDecls(t *testing.T) {
	node := NewQuery()
	node.DeclaredVarDecls = VarDecls{
		NewVarDecl(0, "foo", "string"),
		NewVarDecl(0, "bar", "integer"),
	}
	decls, err := FindVarDecls(node)
	assert.NoError(t, err)
	if assert.Equal(t, len(decls), 5) {
		assert.Equal(t, decls[0].Name, "@eos")
		assert.Equal(t, decls[1].Name, "@eof")
		assert.Equal(t, decls[2].Name, "@timestamp")
		assert.Equal(t, decls[3].Name, "foo")
		assert.Equal(t, decls[4].Name, "bar")
	}
}

// Ensure that duplicate variables in the same scope return an error.
func TestFindVarDeclsDuplicateInScope(t *testing.T) {
	q := NewQuery()
	q.DeclaredVarDecls = VarDecls{
		NewVarDecl(0, "foo", "string"),
		NewVarDecl(0, "foo", "integer"),
	}
	_, err := FindVarDecls(q)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "duplicate symbol in scope: foo")
	}
}
