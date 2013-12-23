package validator

import (
	"github.com/skydb/sky/query/ast"
)

func (v *validator) visitAssignment(n *ast.Assignment, tbl *ast.Symtable) {
	if n.Target == nil {
		v.err = errorf(n, "assignment: target required")
	} else if n.Expression == nil {
		v.err = errorf(n, "assignment: expression required")
	}
}

func (v *validator) exitingAssignment(n *ast.Assignment, tbl *ast.Symtable) {
	targetType := v.dataTypes[n.Target]
	expressionType := v.dataTypes[n.Expression]
	if targetType != expressionType {
		v.err = errorf(n, "assignment: data type mismatch: %s != %s", targetType, expressionType)
	}
}
