package reducer

import (
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hashmap"
)

func (r *Reducer) reduceQuery(q *ast.Query, h *hashmap.Hashmap) error {
	tbl := ast.NewSymtable(nil)
	if err := tbl.Add(q.SystemVarDecls...); err != nil {
		return err
	}
	if err := tbl.Add(q.DeclaredVarDecls...); err != nil {
		return err
	}

	if err := r.reduceStatements(q.Statements, h, tbl); err != nil {
		return err
	}

	return nil
}
