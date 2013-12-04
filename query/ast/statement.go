package ast

type Statement interface {
	Node
	statement()
	String() string
}
