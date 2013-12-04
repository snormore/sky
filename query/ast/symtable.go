package ast

import (
	"fmt"
)

// Symtable represents a scope inside the query and allows variable
// declarations to be looked up while performing codegen.
type Symtable struct {
	parent *Symtable
	locals  map[string]*VarDecl
}

// New creates a new Symtable instance that is associated with the given parent.
// A nil parent means the symtable is the top-level scope.
func NewSymtable(parent *Symtable) *Symtable {
	return &Symtable{
		parent: parent,
		locals:  make(map[string]*VarDecl),
	}
}

// NodeSymtable returns an appropriate symtable for the node.
// If the parent symtable is nil or the node is a block then a new symtable is created.
// Otherwise, the existing symtable is returned.
func NodeSymtable(n Node, parent *Symtable) *Symtable {
	if _, ok := n.(Block); ok || parent == nil {
		return NewSymtable(parent)
	}
	return parent
}

// Locals retrieves a slice of local declarations.
func (tbl *Symtable) Locals() VarDecls {
	locals := make(VarDecls, 0)
	for _, local := range tbl.locals {
		locals = append(locals, local)
	}
	return locals
}

// find looks up a declaration by name. If not found in the current scope
// then declaration is searched for higher up the scope hierarchy.
func (tbl *Symtable) Find(name string) *VarDecl {
	if tbl.locals[name] != nil {
		return tbl.locals[name]
	} else if tbl.parent != nil {
		return tbl.parent.Find(name)
	}
	return nil
}

// add creates a new entry for the declaration in the symbol table. If an
// entry already exists then an error is returned.
func (tbl *Symtable) Add(decl *VarDecl) error {
	if tbl.locals[decl.Name] != nil {
		return fmt.Errorf("duplicate symbol in scope: %s", decl.Name)
	}
	tbl.locals[decl.Name] = decl
	return nil
}
