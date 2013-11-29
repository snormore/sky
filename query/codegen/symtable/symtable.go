package symtable

import (
	"fmt"

	"github.com/skydb/sky/query/ast"
)

// Symtable represents a scope inside the query and allows variable
// declarations to be looked up while performing codegen.
type Symtable struct {
	parent *Symtable
	decls  map[string]*ast.VarDecl
}

// New creates a new Symtable instance that is associated with the given parent.
// A nil parent means the symtable is the top-level scope.
func New(parent *Symtable) *Symtable {
	return &Symtable{
		parent: parent,
		decls:  make(map[string]*ast.VarDecl),
	}
}

// find looks up a declaration by name. If not found in the current scope
// then declaration is searched for higher up the scope hierarchy.
func (tbl *Symtable) Find(name string) *ast.VarDecl {
	if tbl.decls[name] != nil {
		return tbl.decls[name]
	} else if tbl.parent != nil {
		return tbl.parent.Find(name)
	}
	return nil
}

// add creates a new entry for the declaration in the symbol table. If an
// entry already exists then an error is returned.
func (tbl *Symtable) Add(decl *ast.VarDecl) error {
	if tbl.decls[decl.Name] != nil {
		return fmt.Errorf("duplicate symbol in scope: %s", decl.Name)
	}
	tbl.decls[decl.Name] = decl
	return nil
}
