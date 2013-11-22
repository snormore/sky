package ast

type Expression interface {
	QueryElement
	Codegen() (string, error)
	String() string
}
