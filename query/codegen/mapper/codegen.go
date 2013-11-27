package mapper

import (
	"errors"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

var nilValue llvm.Value

// codegen generates LLVM code for a given AST node.
func (m *Mapper) codegen(node ast.Node, tbl *symtable) (llvm.Value, error) {
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

func (m *Mapper) codegenQuery(q *ast.Query, tbl *symtable) (llvm.Value, error) {
	fn := llvm.AddFunction(m.module, "entry", llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false))
	fn.SetFunctionCallConv(llvm.CCallConv)
	entry := llvm.AddBasicBlock(fn, "entry")
	m.builder.SetInsertPointAtEnd(entry)
	m.builder.CreateRet(llvm.ConstInt(llvm.Int32Type(), 12, false))
	return fn, nil
}

