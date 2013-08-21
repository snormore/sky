package query

type BlockStatement interface {
	Statement
	Statements() Statements
}
