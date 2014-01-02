package reducer

import (
	"fmt"

	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hashmap"
)

func (r *Reducer) reduceField(node *ast.Field, h *hashmap.Hashmap, output map[string]interface{}, tbl *ast.Symtable) error {
	identifier := node.Identifier()
	valueType := h.ValueType(hashmap.String(node.Identifier()))

	switch valueType {
	case hashmap.IntValueType:
		prevValue, _ := output[identifier].(int)
		switch node.Aggregation {
		case "count", "sum":
			output[identifier] = prevValue + int(h.Get(hashmap.String(node.Identifier())))
		default:
			return fmt.Errorf("reduce: unsupported int aggregation type: %s", node.Aggregation)
		}

	case hashmap.DoubleValueType:
		prevValue, _ := output[identifier].(float64)
		switch node.Aggregation {
		case "count", "sum":
			output[identifier] = prevValue + h.GetDouble(hashmap.String(node.Identifier()))
		default:
			return fmt.Errorf("reduce: unsupported int aggregation type: %s", node.Aggregation)
		}
	}

	return nil
}
