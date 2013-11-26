package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure that basic variable declarations are found.
func TestFindVarDecls(t *testing.T) {
	node := NewQuery()
	node.DeclaredVarDecls = VarDecls{
		NewVarDecl("foo", "string"),
		NewVarDecl("bar", "integer"),
	}
	decls, err := FindVarDecls(node)
	assert.NoError(t, err)
	if assert.Equal(t, len(decls), 4) {
		assert.Equal(t, decls[0].Name, "@eos")
		assert.Equal(t, decls[1].Name, "@eof")
		assert.Equal(t, decls[2].Name, "foo")
		assert.Equal(t, decls[3].Name, "bar")
	}
}

// Ensure that mismatched variable declaration ids return an error.
func TestFindVarDeclsMismatchId(t *testing.T) {
	node := NewQuery()
	node.DeclaredVarDecls = VarDecls{
		&VarDecl{Id:100, Name:"foo", DataType:"string"},
		&VarDecl{Id:101, Name:"foo", DataType:"string"},
	}
	_, err := FindVarDecls(node)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "Declaration error on 'foo': mismatched id: 100 != 101")
	}
}

// Ensure that mismatched variable declaration data types return an error.
func TestFindVarDeclsMismatchDataType(t *testing.T) {
	node := NewQuery()
	node.DeclaredVarDecls = VarDecls{
		&VarDecl{Name:"foo", DataType:"string"},
		&VarDecl{Name:"foo", DataType:"integer"},
	}
	_, err := FindVarDecls(node)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "Declaration error on 'foo': mismatched data type: string != integer")
	}
}

// Ensure that mismatched variable declaration associations return an error.
func TestFindVarDeclsMismatchAssociation(t *testing.T) {
	node := NewQuery()
	node.DeclaredVarDecls = VarDecls{
		&VarDecl{Name:"foo", DataType:"string", Association:"bar"},
		&VarDecl{Name:"foo", DataType:"string", Association:"baz"},
	}
	_, err := FindVarDecls(node)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), "Declaration error on 'foo': mismatched association: bar != baz")
	}
}
