package reducer

import (
	"fmt"

	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hashmap"
)

func (r *Reducer) reduceField(node *ast.Field, h *hashmap.Hashmap, output map[string]interface{}, tbl *ast.Symtable) error {
	identifier := node.Identifier()
	prevValue, _ := output[identifier].(int)

	switch node.Aggregation {
	case "count", "sum":
		output[identifier] = prevValue + int(h.Get(hashmap.String(node.Identifier())))
	default:
		return fmt.Errorf("reduce: unsupported aggregation type: %s", node.Aggregation)
	}

	return nil
}
