package mapper

import (
	"errors"
	"fmt"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

// codegen generates LLVM code for a given AST node.
func (m *Mapper) codegenStatement(node ast.Statement, symtable *ast.Symtable) (llvm.Value, error) {
	if node == nil {
		return nilValue, errors.New("mapper codegen: unexpected null node")
	}

	symtable = ast.NodeSymtable(node, symtable)

	switch node := node.(type) {
	case *ast.Condition:
		return m.codegenCondition(node, symtable)
	case *ast.EventLoop:
		return m.codegenEventLoop(node, symtable)
	case *ast.Selection:
		return m.codegenSelection(node, symtable)
	case *ast.VarDecl:
		if err := symtable.Add(node); err != nil {
			return nilValue, err
		}
		// return m.codegenVarDecl(node, symtable)
		return nilValue, nil
	default:
		panic(fmt.Sprintf("mapper codegen: unexpected node type: %v", node))
	}
}
