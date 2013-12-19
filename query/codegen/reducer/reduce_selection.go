package reducer

import (
	"fmt"

	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hashmap"
)

func (r *Reducer) reduceSelection(node *ast.Selection, h *hashmap.Hashmap, tbl *ast.Symtable) error {
	output := r.output

	// Drill into name.
	if node.Name != "" {
		h = h.Submap(hashmap.String(node.Name))
		if _, ok := output[node.Name].(map[string]interface{}); !ok {
			output[node.Name] = make(map[string]interface{})
		}
		output = output[node.Name]
	}

	return r.reduceSelectionDimensions(node, h, output, node.Dimensions, tbl)
}

func (r *Reducer) reduceSelectionDimensions(node *ast.Selection, h *hashmap.Hashmap, output map[string]interface{}, dimensions []string, tbl *ast.Symtable) error {
	// Reduce fields if we've reached the end of the dimensions.
	if len(dimensions) == 0 {
		for _, field := range node.Fields {
			if err := r.reduceField(field, h, output, tbl); err != nil {
				return err
			}
		}
		return nil
	}

	// Drill into first dimension
	dimension := dimesions[0]
	decl := tbl.Find(dimension)
	if decl == nil {
		return fmt.Errorf("reduce: dimension not found: %s", dimension)
	}

	// Drill into dimension name.
	h = h.Submap(hashmap.String(dimension))
	if _, ok := output[dimension].(map[string]interface{}); !ok {
		output[dimension] = make(map[string]interface{})
	}
	output = output[node.Name]

	// Iterate over dimension values.
	iterator := hashmap.NewIterator(h)
	for {
		key, ok := iterator.Next()
		if !ok {
			break
		}

		// Convert value to appropriate type based on variable decl.
		var keyString string
		switch decl.DataType {
		case core.StringDataType:
			return fmt.Errorf("reduce: string dimensions are not supported: %s", dimension)
		case core.FloatDataType:
			return fmt.Errorf("reduce: float dimensions are not supported: %s", dimension)
		case core.IntegerDataType:
			keyString = strconv.Itoa(key)
		case core.FactorDataType:
			// TODO: Defactorize!
			keyString = strconv.Itoa(key)
		case core.BooleanDataType:
			if key == 0 {
				keyString = "false"
			} else {
				keyString = "true"
			}
		}

		// Drill into output map.
		if _, ok := output[keyString].(map[string]interface{}); !ok {
			output[keyString] = make(map[string]interface{})
		}

		// Recursively drill into next dimension.
		if err := m.reduceSelectionDimensions(node, h.Submap(key), output[keyString], dimensions[1:], tbl); err != nil {
			return err
		}
	}

	return nil
}
