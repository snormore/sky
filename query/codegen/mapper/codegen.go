package mapper

import (
	"errors"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

var nilValue llvm.Value

// codegen generates LLVM code for a given AST node.
func (m *Mapper) codegen(node ast.Node, symtable *ast.Symtable) (llvm.Value, error) {
	if node == nil {
		return nilValue, errors.New("mapper codegen: unexpected null node")
	}

	symtable = ast.NodeSymtable(node, symtable)

	switch node := node.(type) {
	case *ast.Query:
		return m.codegenQuery(node, symtable)
	case *ast.VarDecl:
		if err := symtable.Add(node); err != nil {
			return nilValue, err
		}
		// return m.codegenVarDecl(node, symtable)
		return nilValue, nil
	default:
		panic("mapper codegen: unexpected node type")
	}
}
