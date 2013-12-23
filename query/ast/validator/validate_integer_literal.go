package validator

import (
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

func (v *validator) visitIntegerLiteral(n *ast.IntegerLiteral, tbl *ast.Symtable) {
	v.dataTypes[n] = core.IntegerDataType
}
