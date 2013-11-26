package mapper

import (
	"errors"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

// codegen generates LLVM code for a given AST node.
func codegen(m *llvm.Module, tbl *symtable, node ast.Node) error {
	if node == nil {
		return errors.New("mapper codegen: unexpected null node")
	}

	switch node := node.(type) {
	case *ast.Query:
		return codegenQuery(m, tbl, node)
	default:
		panic("mapper codegen: unexpected node type:")
	}
}

func codegenQuery(m *llvm.Module, tbl *symtable, q *ast.Query) error {
	builder := llvm.NewBuilder()
	defer builder.Dispose()
	return nil
}
