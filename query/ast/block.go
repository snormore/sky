package ast

// Block represents an AST node that creates a new variable scope.
type Block interface {
	Node
	block()
}
