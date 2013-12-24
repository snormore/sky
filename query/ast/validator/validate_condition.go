package validator

import (
	"github.com/skydb/sky/query/ast"
)

func (v *validator) visitCondition(n *ast.Condition, tbl *ast.Symtable) {
	switch n.UOM {
	case ast.UnitSteps, "":
	default:
		v.err = errorf(n, "condition: unit not supported: %s", n.UOM)
		return
	}

	if n.Start > n.End {
		v.err = errorf(n, "condition: start index (%d) cannot be less than end index (%d)", n.Start, n.End)
	}
}

func (v *validator) exitingCondition(n *ast.Condition, tbl *ast.Symtable) {
	if n.Expression != nil {
		expressionType := v.dataTypes[n.Expression]
		if expressionType != "boolean" {
			v.err = errorf(n, "condition: invalid expression return type: %s", expressionType)
		}
	}
}
