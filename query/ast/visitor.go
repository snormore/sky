package ast

// Visitor is the interface that is used for walking the AST tree.
type Visitor interface {
	Visit(node Node) Visitor
}
