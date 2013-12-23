package validator

import (
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

func (v *validator) visitStringLiteral(n *ast.StringLiteral, tbl *ast.Symtable) {
	v.dataTypes[n] = core.FactorDataType
}
