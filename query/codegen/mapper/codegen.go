package mapper

import (
	"errors"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/symtable"
)

var nilValue llvm.Value

// codegen generates LLVM code for a given AST node.
func (m *Mapper) codegen(node ast.Node, tbl *symtable.Symtable) (llvm.Value, error) {
	if node == nil {
		return nilValue, errors.New("mapper codegen: unexpected null node")
	}

	switch node := node.(type) {
	case *ast.Query:
		return m.codegenQuery(node, tbl)
	default:
		panic("mapper codegen: unexpected node type")
	}
}

