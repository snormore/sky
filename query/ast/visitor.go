package ast

// Visitor is the interface that is used for walking the AST tree.
type Visitor interface {
	Visit(Node, *Symtable) Visitor
	Error(error) Visitor
}

// ExitingVisitor is the interface that is used for walking the AST tree
// but only receiving callbacks after the node has been processed.
type ExitingVisitor interface {
	Visitor
	Exiting(Node, *Symtable)
}
