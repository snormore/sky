package ast

import (
	"sort"
	"strings"

	"github.com/skydb/sky/core"
)

// Query represents a collection of statements used to process data in
// the Sky database.
type Query struct {
	Prefix           string
	SystemVarDecls   VarDecls
	DeclaredVarDecls VarDecls
	Statements       Statements
}

func (q *Query) node() {}
func (q *Query) block() {}

// NewQuery returns a new query.
func NewQuery() *Query {
	q := &Query{}
	q.SystemVarDecls = VarDecls{
		NewVarDecl(0, "@eos", core.BooleanDataType),
		NewVarDecl(0, "@eof", core.BooleanDataType),
		NewVarDecl(0, "@timestamp", core.IntegerDataType),
	}
	return q
}

// VarDecls returns a list of indexed variable declarations.
func (q *Query) VarDecls() (VarDecls, error) {
	// Retrieve list of 
	decls, err := FindVarDecls(q)
	if err != nil {
		return nil, err
	}

	// Sort and index declarations.
	sort.Sort(decls)
	for i, decl := range decls {
		decl.index = i
	}

	return decls, nil
}

func (q *Query) String() string {
	arr := []string{}
	for _, v := range q.DeclaredVarDecls {
		arr = append(arr, v.String())
	}
	arr = append(arr, q.Statements.String())
	return strings.Join(arr, "\n")
}
