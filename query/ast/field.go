package ast

import (
	"fmt"
	"strings"
)

// Field represents a single field within a Selection.
type Field struct {
	Name        string
	Aggregation string
	Distinct    bool
	Expression  Expression
}

func (f *Field) node() string {}

// NewField creates a new Field instance.
func NewField(name string, aggregation string, expression Expression) *Field {
	return &Field{
		Name:        name,
		Aggregation: aggregation,
		Expression:  expression,
	}
}

// Determines if the field performs aggregation.
func (f *Field) IsAggregate() bool {
	return f.Aggregation != ""
}

//--------------------------------------
// Finalization
//--------------------------------------

// Finalizes the results into a final state after merge.
func (f *Field) Finalize(data interface{}) error {
	switch f.Aggregation {
	case "count":
		if f.Distinct {
			aggregation := data.(map[interface{}]interface{})
			object := aggregation[f.Name].(map[interface{}]interface{})
			value := object["distinct"]
			aggregation[f.Name] = value
		}
	case "avg":
		aggregation := data.(map[interface{}]interface{})
		object := aggregation[f.Name].(map[interface{}]interface{})
		value := object["avg"]
		aggregation[f.Name] = value
	}
	return nil
}

// Converts the field to a string-based representation.
func (f *Field) String() string {
	var expr string
	if f.expression != nil {
		expr = f.expression.String()
	}

	var str string
	if f.Aggregation == "" {
		str = expr
	} else {
		if f.Distinct {
			str = fmt.Sprintf("%s(DISTINCT %s)", f.Aggregation, expr)
		} else {
			str = fmt.Sprintf("%s(%s)", f.Aggregation, expr)
		}
	}
	if f.Name != "" {
		str += " AS " + f.Name
	}
	return str
}
