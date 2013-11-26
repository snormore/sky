package ast

import (
	"fmt"
)

// Field represents a single field within a Selection.
type Field struct {
	Name        string
	Aggregation string
	Distinct    bool
	Expression  Expression
}

func (f *Field) node() {}

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

// Converts the field to a string-based representation.
func (f *Field) String() string {
	var expr string
	if f.Expression != nil {
		expr = f.Expression.String()
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
