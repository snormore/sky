package query

type Expression interface {
	QueryElement
	Codegen() (string, error)
	String() string
}
