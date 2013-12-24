package ast

// ForEach executes a function over every node in the AST tree.
func ForEach(node Node, fn func(Node)) {
	v := new(forEachVisitor)
	v.fn = fn
	Walk(v, node)
}

// forEachVisitor is an AST visitor that executes a function over every AST node.
type forEachVisitor struct {
	fn func(Node)
}

func (v *forEachVisitor) Visit(node Node, symtable *Symtable) Visitor {
	v.fn(node)
	return v
}

func (v *forEachVisitor) Error(err error) Visitor {
	return v
}
