package validator

import (
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

func (v *validator) exitingBinaryExpression(n *ast.BinaryExpression, tbl *ast.Symtable) {
	lhsType := v.dataTypes[n.LHS]
	rhsType := v.dataTypes[n.RHS]

	// Raise error for type mismatch (except for "float-integer").
	if lhsType != rhsType {
		if !((lhsType == "float" && rhsType == "integer") || (lhsType == "integer" && rhsType == "float")) {
			v.err = errorf(n, "expression: data type mismatch: %s != %s", lhsType, rhsType)
		}
	}

	switch lhsType {
	case core.BooleanDataType:
		v.exitingBooleanBinaryExpression(n, tbl)
	case core.FactorDataType:
		v.exitingFactorBinaryExpression(n, tbl)
	case core.IntegerDataType, core.FloatDataType:
		v.exitingNumericBinaryExpression(n, tbl)
	default:
		v.err = errorf(n, "expression: invalid binary expression type: %s", lhsType)
	}
}

func (v *validator) exitingBooleanBinaryExpression(n *ast.BinaryExpression, tbl *ast.Symtable) {
	switch n.Op {
	case ast.OpEquals, ast.OpNotEquals, ast.OpAnd, ast.OpOr:
		v.dataTypes[n] = core.BooleanDataType
	default:
		v.err = errorf(n, "expression: invalid boolean operator: %s", n.OpString())
	}
}

func (v *validator) exitingFactorBinaryExpression(n *ast.BinaryExpression, tbl *ast.Symtable) {
	lhsVarRef, _ := n.LHS.(*ast.VarRef)
	lhsString, _ := n.LHS.(*ast.StringLiteral)
	rhsVarRef, _ := n.RHS.(*ast.VarRef)
	rhsString, _ := n.RHS.(*ast.StringLiteral)

	// Check that two refs associated with each other point to the same type of factor.
	if lhsVarRef != nil && rhsVarRef != nil {
		lhsDecl := tbl.Find(lhsVarRef.Name)
		rhsDecl := tbl.Find(rhsVarRef.Name)
		if lhsDecl.Name != rhsDecl.Name && lhsDecl.Name != rhsDecl.Association && lhsDecl.Association != rhsDecl.Name {
			v.err = errorf(n, "expression: mismatched factor association: %s<%s> != %s<%s>", lhsDecl.Name, lhsDecl.Association, rhsDecl.Name, rhsDecl.Association)
			return
		}
	}

	// Make sure two strings are not compared.
	if lhsString != nil && rhsString != nil {
		v.err = errorf(n, "expression: string literal comparison not allowed: %s", n.String())
		return
	}

	switch n.Op {
	case ast.OpEquals, ast.OpNotEquals:
		v.dataTypes[n] = core.BooleanDataType
	default:
		v.err = errorf(n, "expression: invalid factor operator: %s", n.OpString())
	}
}

func (v *validator) exitingNumericBinaryExpression(n *ast.BinaryExpression, tbl *ast.Symtable) {
	switch n.Op {
	case ast.OpEquals, ast.OpNotEquals:
		v.dataTypes[n] = core.BooleanDataType
	case ast.OpGreaterThan, ast.OpGreaterThanOrEqualTo:
		v.dataTypes[n] = core.BooleanDataType
	case ast.OpLessThan, ast.OpLessThanOrEqualTo:
		v.dataTypes[n] = core.BooleanDataType
	case ast.OpPlus, ast.OpMinus:
		v.dataTypes[n] = core.IntegerDataType
	case ast.OpMultiply, ast.OpDivide:
		v.dataTypes[n] = core.IntegerDataType
	default:
		v.err = errorf(n, "expression: invalid numeric operator: %s", n.OpString())
	}
}
