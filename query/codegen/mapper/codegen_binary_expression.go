package mapper

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

func (m *Mapper) codegenBinaryExpression(node *ast.BinaryExpression, event llvm.Value, tbl *ast.Symtable) (llvm.Value, error) {
	// Special handling for factor expressions.
	if m.isFactorVarRef(node.LHS, tbl) {
		return m.codegenFactorEquality(node, node.LHS.(*ast.VarRef), node.RHS, event, tbl)
	} else if m.isFactorVarRef(node.RHS, tbl) {
		return m.codegenFactorEquality(node, node.RHS.(*ast.VarRef), node.LHS, event, tbl)
	}

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

func (m *Mapper) isFactorVarRef(node ast.Expression, tbl *ast.Symtable) bool {
	ref, ok := node.(*ast.VarRef)
	if !ok {
		return false
	}
	decl := tbl.Find(ref.Name)
	return (decl != nil && decl.DataType == core.FactorDataType)
}

func (m *Mapper) codegenFactorEquality(node *ast.BinaryExpression, lhs *ast.VarRef, rhs ast.Expression, event llvm.Value, tbl *ast.Symtable) (llvm.Value, error) {
	// Only allow == for factors.
	if node.Op != ast.OpEquals {
		return nilValue, fmt.Errorf("Non-equality operators not allowed for factor data types: %s", node.String())
	}

	lhsValue, err := m.codegenExpression(lhs, event, tbl)
	if err != nil {
		return nilValue, err
	}

	// Only allow comparisons to strings or other associated factors.
	switch rhs := rhs.(type) {
	case *ast.StringLiteral:
		id, err := m.factorizer.Factorize(lhs.Name, rhs.Value, false)
		if err != nil {
			return nilValue, err
		}
		return m.icmp(llvm.IntEQ, lhsValue, m.constint(int(id)), ""), nil

	case *ast.VarRef:
		rhsValue, err := m.codegenExpression(rhs, event, tbl)
		if err != nil {
			return nilValue, err
		}
		return m.icmp(llvm.IntEQ, lhsValue, rhsValue, ""), nil

	default:
		return nilValue, fmt.Errorf("Invalid factor expression: %s", node.String())
	}
}
