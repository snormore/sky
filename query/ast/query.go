package ast

import (
	"bytes"
	"fmt"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/db"
	"sort"
	"strings"
)

// Query represents a collection of statements used to process data in
// the Sky database.
type Query struct {
	Factorizer        db.Factorizer
	Table             *core.Table
	Prefix            string
	SystemVariables   []*Variable
	DeclaredVariables []*Variable
	Statements        Statements
}

func (q *Query) node() string {}

// NewQuery returns a new query.
func NewQuery() *Query {
	q := &Query{}
	q.SystemVariables = []*Variable{
		NewVariable("@eos", core.BooleanDataType),
		NewVariable("@eof", core.BooleanDataType),
	}
	return q
}

func (q *Query) String() string {
	arr := []string{}
	for _, v := range q.DeclaredVariables {
		arr = append(arr, v.String())
	}
	arr = append(arr, q.Statements.String())
	return strings.Join(arr, "\n")
}
