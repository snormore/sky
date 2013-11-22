package ast

type BlockStatement interface {
	Statement
	Statements() Statements
}
