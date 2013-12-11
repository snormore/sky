package mapper

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

func (m *Mapper) codegenBinaryExpression(node *ast.BinaryExpression, event llvm.Value, tbl *ast.Symtable) (llvm.Value, error) {
	lhs, err := m.codegenExpression(node.LHS, event, tbl)
	if err != nil {
		return nilValue, err
	}

	rhs, err := m.codegenExpression(node.RHS, event, tbl)
	if err != nil {
		return nilValue, err
	}

	switch node.Op {
	case ast.OpEquals:
		return m.icmp(llvm.IntEQ, lhs, rhs, ""), nil
	case ast.OpNotEquals:
		return m.icmp(llvm.IntNE, lhs, rhs, ""), nil
	case ast.OpGreaterThan:
		return m.icmp(llvm.IntSGT, lhs, rhs, ""), nil
	case ast.OpGreaterThanOrEqualTo:
		return m.icmp(llvm.IntSGE, lhs, rhs, ""), nil
	case ast.OpLessThan:
		return m.icmp(llvm.IntSLT, lhs, rhs, ""), nil
	case ast.OpLessThanOrEqualTo:
		return m.icmp(llvm.IntSLE, lhs, rhs, ""), nil
	case ast.OpPlus:
		return m.builder.CreateAdd(lhs, rhs, ""), nil
	case ast.OpMinus:
		return m.builder.CreateSub(lhs, rhs, ""), nil
	case ast.OpMultiply:
		return m.builder.CreateMul(lhs, rhs, ""), nil
	case ast.OpDivide:
		return m.builder.CreateSDiv(lhs, rhs, ""), nil
	case ast.OpAnd:
		return m.builder.CreateAnd(lhs, rhs, ""), nil
	case ast.OpOr:
		return m.builder.CreateOr(lhs, rhs, ""), nil
	default:
	}

	return nilValue, fmt.Errorf("Invalid binary expression operator: %d", node.Op)
}
