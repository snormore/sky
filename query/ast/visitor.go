package ast

type Visitor interface {
	Visit(node Node) (w Visitor)
}
