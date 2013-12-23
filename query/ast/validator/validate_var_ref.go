package validator

import (
	"github.com/skydb/sky/query/ast"
)

func (v *validator) visitVarRef(n *ast.VarRef, tbl *ast.Symtable) {
	if n.Name == "" {
		v.err = errorf(n, "variable: name required")
	} else if decl := tbl.Find(n.Name); decl == nil {
		v.err = errorf(n, "variable: not found: %s", n.Name)
	} else {
		v.dataTypes[n] = decl.DataType
	}
}
