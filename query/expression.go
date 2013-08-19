package query

type Expression interface {
	QueryElement
	String() string
}
