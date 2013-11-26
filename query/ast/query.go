package ast

import (
	"strings"

	"github.com/skydb/sky/core"
	"github.com/skydb/sky/db"
)

// Query represents a collection of statements used to process data in
// the Sky database.
type Query struct {
	Factorizer       db.Factorizer
	Table            *core.Table
	Prefix           string
	SystemVarDecls   VarDecls
	DeclaredVarDecls VarDecls
	Statements       Statements
}

func (q *Query) node() {}

// NewQuery returns a new query.
func NewQuery() *Query {
	q := &Query{}
	q.SystemVarDecls = VarDecls{
		NewVarDecl("@eos", core.BooleanDataType),
		NewVarDecl("@eof", core.BooleanDataType),
	}
	return q
}

func (q *Query) String() string {
	arr := []string{}
	for _, v := range q.DeclaredVarDecls {
		arr = append(arr, v.String())
	}
	arr = append(arr, q.Statements.String())
	return strings.Join(arr, "\n")
}
