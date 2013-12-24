package validator

import (
	"github.com/skydb/sky/query/ast"
)

func (v *validator) visitQuery(n *ast.Query, tbl *ast.Symtable) {
	selections := n.Selections()

	// Validate that all selections use consistent naming.
	var named interface{}
	for _, selection := range selections {
		if (selection.Name == "" && named == true) || (selection.Name != "" && named == false) {
			v.err = errorf(n, "query: mixing named and unnamed selections is not allowed")
			return
		}
		named = (selection.Name != "")
	}

	// Validate all selections use consistent aggregate/non-aggregate.
	// NOTE: Multiple aggregation types within a selection is checked at the selection.
	var aggregated interface{}
	for _, selection := range selections {
		if (selection.HasAggregateFields() && aggregated == false) || (selection.HasNonAggregateFields() && aggregated == true) {
			v.err = errorf(n, "query: mixing aggregated and non-aggregated selection is not allowed")
			return
		}
		aggregated = selection.HasAggregateFields()
	}
}

func (v *validator) exitingQuery(n *ast.Query, tbl *ast.Symtable) {
}
