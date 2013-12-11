package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

func (m *Mapper) codegenBooleanLiteral(node *ast.BooleanLiteral, event llvm.Value, tbl *ast.Symtable) (llvm.Value, error) {
	return m.constbool(node.Value), nil
}
