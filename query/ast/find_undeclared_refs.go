package ast

// FindUndeclaredRefs retrieves a list of referenced variables without declarations.
func FindUndeclaredRefs(node Node) []*VarRef {
	v := &undeclaredRefsVisitor{make(map[string]*VarRef)}
	Walk(v, node)

	refs := make([]*VarRef, 0)
	for _, ref := range v.refs {
		refs = append(refs, ref)
	}
	return refs
}

// undeclaredRefsVisitor is an AST visitor that finds variable references without declarations.
type undeclaredRefsVisitor struct {
	refs map[string]*VarRef
}

func (v *undeclaredRefsVisitor) Visit(node Node, symtable *Symtable) Visitor {
	if node, ok := node.(*VarRef); ok {
		name := node.Name
		if symtable.Find(name) == nil && v.refs[name] == nil {
			v.refs[name] = node
		}
	}
	return v
}

func (v *undeclaredRefsVisitor) Error(err error) Visitor {
	return v
}
