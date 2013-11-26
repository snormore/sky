package ast

// Expression represents an AST node that evaluates to a value.
type Expression interface {
	Node
	expression()
	String() string
}
