package skyd

import (
	"fmt"
)

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	QueryStepTypeCondition = "condition"
	QueryStepTypeSelection = "selection"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

type QueryStep interface {
	FunctionName() string
	MergeFunctionName() string
	Serialize() map[string]interface{}
	Deserialize(map[string]interface{}) error
	Codegen() (string, error)
}

type QueryStepList []QueryStep

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query step list into an untyped object.
func (l QueryStepList) Serialize() []interface{} {
	steps := make([]interface{}, 0)
	for _, step := range l {
		steps = append(steps, step.Serialize())
	}
	return steps
}

// Decodes a query step list from an untyped slice.
func DeserializeQueryStepList(obj interface{}, q *Query) (QueryStepList, error) {
	l := make(QueryStepList, 0)
	if steps, ok := obj.([]interface{}); ok {
		for _, _s := range steps {
			if s, ok := _s.(map[string]interface{}); ok {
				var step QueryStep
				switch s["type"] {
				case QueryStepTypeCondition:
					step = NewQueryCondition(q)
				case QueryStepTypeSelection:
					step = NewQuerySelection(q)
				default:
					return nil, fmt.Errorf("Invalid query step type: %v", s["type"])
				}
				err := step.Deserialize(s)
				if err != nil {
					return nil, err
				}
				l = append(l, step)
			} else {
				return nil, fmt.Errorf("Invalid step: %v", obj)
			}
		}
	} else {
		return nil, fmt.Errorf("Invalid steps: %v", obj)
	}
	return l, nil
}
