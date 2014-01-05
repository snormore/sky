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
	reducible   bool
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

// Reducible returns whether the field should be reduced. This can be false if
// a field with the same path has already been declared previously.
func (f *Field) Reducible() bool {
	return f.reducible
}

// Returns the identifier used by the field.
// This is automatically generated if a name is not explicitly set.
func (f *Field) Identifier() string {
	if f.Name != "" {
		return f.Name
	}
	names := []string{}
	if f.Aggregation != "" {
		names = append(names, f.Aggregation)
	}
	if f.Expression != nil {
		expr := f.Expression.String()
		expr = nonAlphaRegex.ReplaceAllString(expr, "_")
		expr = strings.Trim(expr, "_")
		names = append(names, expr)
	}
	return strings.Join(names, "_")
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
