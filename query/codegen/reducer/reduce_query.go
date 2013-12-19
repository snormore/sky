package reducer

import (
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hashmap"
)

func (r *Reducer) reduceQuery(q *ast.Query, h *hashmap.Hashmap) error {
	tbl := ast.NewSymtable(nil)

	if err := r.reduceStatements(q.Statements, h, tbl); err != nil {
		return err
	}

	return nil
}
