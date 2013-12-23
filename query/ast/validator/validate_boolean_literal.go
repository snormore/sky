package validator

import (
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

func (v *validator) visitBooleanLiteral(n *ast.BooleanLiteral, tbl *ast.Symtable) {
	v.dataTypes[n] = core.BooleanDataType
}
