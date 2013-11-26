package ast

type Statement interface {
	Node
	String() string
}
