package ast

// Find retrieves a list of nodes matching a given filter function.
func Find(node Node, fn func(Node) bool) []Node {
	v := new(findVisitor)
	v.fn = fn
	Walk(v, node)
	return v.nodes
}

// findVisitor is an AST visitor that collects a list of nodes matching a filter function.
type findVisitor struct {
	nodes []Node
	fn    func(Node) bool
}

func (v *findVisitor) Visit(node Node, symtable *Symtable) Visitor {
	if v.fn(node) {
		v.nodes = append(v.nodes, node)
	}
	return v
}

func (v *findVisitor) Error(err error) Visitor {
	return v
}
