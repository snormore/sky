package query

import (
	"errors"
	"fmt"
	"regexp"
)

type SelectionField struct {
	queryElementImpl
	Name        string
	Aggregation string
	expression  Expression
}

// Creates a new selection field.
func NewSelectionField(name string, aggregation string, expression Expression) *SelectionField {
	f := &SelectionField{
		Name:        name,
		Aggregation: aggregation,
	}
	f.SetExpression(expression)
	return f
}

// Returns the expression evaluated by the field.
func (f *SelectionField) Expression() Expression {
	return f.expression
}

// Sets the expression.
func (f *SelectionField) SetExpression(expression Expression) {
	if f.expression != nil {
		f.expression.SetParent(nil)
	}
	f.expression = expression
	if f.expression != nil {
		f.expression.SetParent(f)
	}
}

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query selection into an untyped map.
func (f *SelectionField) Serialize() map[string]interface{} {
	obj := map[string]interface{}{
		"name":        f.Name,
		"aggregation": f.Aggregation,
	}
	if f.expression != nil {
		obj["expression"] = f.expression.String()
	}
	return obj
}

// Decodes a query selection from an untyped map.
func (f *SelectionField) Deserialize(obj map[string]interface{}) error {
	if obj == nil {
		return errors.New("SelectionField: Unable to deserialize nil.")
	}

	// Deserialize "aggregation".
	if obj["aggregation"] == nil {
		f.Aggregation = ""
	} else if aggregation, ok := obj["aggregation"].(string); ok {
		f.Aggregation = aggregation
	} else {
		return fmt.Errorf("SelectionField: Invalid aggregation: %v", obj["aggregation"])
	}

	// Deserialize "expression".
	if obj["expression"] == nil {
		f.SetExpression(nil)
	} else if expression, ok := obj["expression"].(string); ok {
		// Extract the aggregation from the expression for backwards compatibility.
		if m := regexp.MustCompile(`^(\w+)\((.*)\)$`).FindStringSubmatch(expression); m != nil && len(f.Aggregation) == 0 {
			f.Aggregation = m[1]
			expression = m[2]
		}

		// Parse expression.
		if len(expression) > 0 {
			expr, err := NewExpressionParser().ParseString(expression)
			if err != nil {
				return err
			}
			f.SetExpression(expr)
		} else {
			f.SetExpression(nil)
		}
	} else {
		return fmt.Errorf("SelectionField: Invalid expression: %v", obj["expression"])
	}

	// Deserialize "name".
	if name, ok := obj["name"].(string); ok && len(name) > 0 {
		f.Name = name
	} else {
		return fmt.Errorf("SelectionField: Invalid name: %v", obj["name"])
	}

	return nil
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the expression.
func (f *SelectionField) CodegenExpression(init bool) (string, error) {
	var code string
	var err error
	if f.expression != nil {
		if code, err = f.expression.Codegen(); err != nil {
			return "", err
		}
	}

	// Expressions are required for everything except count().
	if f.expression == nil && f.Aggregation != "count" {
		return "", fmt.Errorf("Selection field expression required for '%s'", f.Aggregation)
	}

	switch f.Aggregation {
	case "sum":
		return fmt.Sprintf("data.%s = (data.%s or 0) + (%s)", f.Name, f.Name, code), nil
	case "min":
		return fmt.Sprintf("if(data.%s == nil or data.%s > (%s)) then data.%s = (%s) end", f.Name, f.Name, code, f.Name, code), nil
	case "max":
		return fmt.Sprintf("if(data.%s == nil or data.%s < (%s)) then data.%s = (%s) end", f.Name, f.Name, code, f.Name, code), nil
	case "count":
		return fmt.Sprintf("data.%s = (data.%s or 0) + 1", f.Name, f.Name), nil
	case "avg":
		return fmt.Sprintf("if data.%s == nil then data.%s = sky_average_new() end sky_average_insert(data.%s, (%s))", f.Name, f.Name, f.Name, code), nil
	case "histogram":
		if init {
			return fmt.Sprintf("if data.%s == nil then data.%s = sky_histogram_new() end table.insert(data.%s.values, (%s))", f.Name, f.Name, f.Name, code), nil
		} else {
			return fmt.Sprintf("if data.%s ~= nil then sky_histogram_insert(data.%s, (%s)) end", f.Name, f.Name, code), nil
		}
	}

	return "", fmt.Errorf("SelectionField: Unsupported aggregation method: %v", f.Aggregation)
}

// Generates Lua code for the merge expression.
func (f *SelectionField) CodegenMergeExpression() (string, error) {
	switch f.Aggregation {
	case "sum":
		return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", f.Name, f.Name, f.Name), nil
	case "min":
		return fmt.Sprintf("if(result.%s == nil or result.%s > data.%s) then result.%s = data.%s end", f.Name, f.Name, f.Name, f.Name, f.Name), nil
	case "max":
		return fmt.Sprintf("if(result.%s == nil or result.%s < data.%s) then result.%s = data.%s end", f.Name, f.Name, f.Name, f.Name, f.Name), nil
	case "count":
		return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", f.Name, f.Name, f.Name), nil
	case "avg":
		return fmt.Sprintf("result.%s = sky_average_merge(result.%s, data.%s)", f.Name, f.Name, f.Name), nil
	case "histogram":
		return fmt.Sprintf("result.%s = sky_histogram_merge(result.%s, data.%s)", f.Name, f.Name, f.Name), nil
	case "":
		return fmt.Sprintf("result.%s = data.%s", f.Name, f.Name), nil
	}

	return "", fmt.Errorf("SelectionField: Unsupported merge aggregation method: %s", f.Aggregation)
}

//--------------------------------------
// Initialization
//--------------------------------------

// Checks if the field requires the data structure to be initialized before
// aggregation. This will occur when computing histograms since all servlets
// need to insert into the same bins.
func (f *SelectionField) RequiresInitialization() bool {
	return f.Aggregation == "histogram"
}

//--------------------------------------
// Utility
//--------------------------------------

// Retrieves a list of variables referenced by this field.
func (f *SelectionField) VarRefs() []*VarRef {
	refs := []*VarRef{}
	if f.expression != nil {
		refs = append(refs, f.expression.VarRefs()...)
	}
	return refs
}

// Converts the field to a string-based representation.
func (f *SelectionField) String() string {
	var expr string
	if f.expression != nil {
		expr = f.expression.String()
	}
	return fmt.Sprintf("%s(%s) AS %s", f.Aggregation, expr, f.Name)
}
