package skyd

import (
	"errors"
	"fmt"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A selection step aggregates data in a query.
type QuerySelection struct {
	query      *Query
	functionName string
	mergeFunctionName string
	Expression string
	Alias      string
	Dimensions []string
	Steps      QueryStepList
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// Creates a new selection.
func NewQuerySelection(query *Query) *QuerySelection {
	id := query.NextIdentifier()
	return &QuerySelection{
		query: query,
		functionName: fmt.Sprintf("f%d", id),
		mergeFunctionName: fmt.Sprintf("m%d", id),
	}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// Retrieves the query this selection is associated with.
func (s *QuerySelection) Query() *Query {
	return s.query
}

// Retrieves the function name used during codegen.
func (s *QuerySelection) FunctionName() string {
	return s.functionName
}

// Retrieves the merge function name used during codegen.
func (s *QuerySelection) MergeFunctionName() string {
	return s.mergeFunctionName
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query selection into an untyped map.
func (s *QuerySelection) Serialize() map[string]interface{} {
	obj := map[string]interface{}{
		"type":       QueryStepTypeSelection,
		"expression": s.Expression,
		"alias":      s.Alias,
		"dimensions": s.Dimensions,
		"steps":      s.Steps.Serialize(),
	}
	return obj
}

// Decodes a query selection from an untyped map.
func (s *QuerySelection) Deserialize(obj map[string]interface{}) error {
	if obj == nil {
		return errors.New("skyd.QuerySelection: Unable to deserialize nil.")
	}
	if obj["type"] != QueryStepTypeSelection {
		return fmt.Errorf("skyd.QuerySelection: Invalid step type: %v", obj["type"])
	}

	// Deserialize "expression".
	if expression, ok := obj["expression"].(string); ok && len(expression) > 0 {
		s.Expression = expression
	} else {
		return fmt.Errorf("skyd.QuerySelection: Invalid expression: %v", obj["expression"])
	}

	// Deserialize "alias".
	if alias, ok := obj["alias"].(string); ok && len(alias) > 0 {
		s.Alias = alias
	} else {
		return fmt.Errorf("skyd.QuerySelection: Invalid alias: %v", obj["alias"])
	}

	// Deserialize "dimensions".
	if dimensions, ok := obj["dimensions"].([]interface{}); ok {
		s.Dimensions = []string{}
		for _, dimension := range dimensions {
			if str, ok := dimension.(string); ok {
				s.Dimensions = append(s.Dimensions, str)
			} else {
				return fmt.Errorf("skyd.QuerySelection: Invalid dimension: %v", dimension)
			}
		}
	} else {
		return fmt.Errorf("skyd.QuerySelection: Invalid dimensions: %v", obj["dimensions"])
	}

	// Deserialize steps.
	var err error
	s.Steps, err = DeserializeQueryStepList(obj["steps"], s.query)
	if err != nil {
		return err
	}

	return nil
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the query.
func (s *QuerySelection) Codegen() (string, error) {
	return "", nil
}
