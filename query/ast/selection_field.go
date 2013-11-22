package ast

import (
	"fmt"
	"strings"
)

type SelectionField struct {
	queryElementImpl
	Name        string
	Aggregation string
	Distinct    bool
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

// Determines if the field performs aggregation.
func (f *SelectionField) IsAggregate() bool {
	return f.Aggregation != ""
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
// Code Generation
//--------------------------------------

// Returns the name used by the field. This is automatically generated
// if a name is not explicitly set.
func (f *SelectionField) CodegenName() string {
	if f.Name != "" {
		return f.Name
	} else {
		names := []string{}
		if f.Aggregation != "" {
			names = append(names, f.Aggregation)
		}
		if f.expression != nil {
			expr := f.expression.String()
			expr = nonAlphaRegex.ReplaceAllString(expr, "_")
			expr = strings.Trim(expr, "_")
			names = append(names, expr)
		}
		return strings.Join(names, "_")
	}
}

// Generates Lua code for the expression.
func (f *SelectionField) CodegenExpression(init bool) (string, error) {
	var code string
	var err error
	name := f.CodegenName()

	if f.expression != nil {
		if code, err = f.expression.Codegen(); err != nil {
			return "", err
		}
	}

	if f.IsAggregate() {
		// Expressions are required for everything except count().
		if f.expression == nil && f.Aggregation != "count" && !f.Distinct {
			return "", fmt.Errorf("Selection field expression required for '%s'", f.Aggregation)
		}

		switch f.Aggregation {
		case "sum":
			return fmt.Sprintf("data.%s = (data.%s or 0) + (%s)", name, name, code), nil
		case "min":
			return fmt.Sprintf("if(data.%s == nil or data.%s > (%s)) then data.%s = (%s) end", name, name, code, name, code), nil
		case "max":
			return fmt.Sprintf("if(data.%s == nil or data.%s < (%s)) then data.%s = (%s) end", name, name, code, name, code), nil
		case "count":
			if f.Distinct {
				return fmt.Sprintf("if data.%s == nil then data.%s = sky_distinct_new() end sky_distinct_insert(data.%s, (%s))", name, name, name, code), nil
			} else {
				return fmt.Sprintf("data.%s = (data.%s or 0) + 1", name, name), nil
			}
		case "avg":
			return fmt.Sprintf("if data.%s == nil then data.%s = sky_average_new() end sky_average_insert(data.%s, (%s))", name, name, name, code), nil
		case "histogram":
			if init {
				return fmt.Sprintf("if data.%s == nil then data.%s = sky_histogram_new() end table.insert(data.%s.values, (%s))", name, name, name, code), nil
			} else {
				return fmt.Sprintf("if data.%s ~= nil then sky_histogram_insert(data.%s, (%s)) end", name, name, code), nil
			}
		}

		return "", fmt.Errorf("SelectionField: Unsupported aggregation method: %v", f.Aggregation)

	} else {
		return fmt.Sprintf("row.%s = %s", name, code), nil
	}
}

// Generates Lua code for the merge expression.
func (f *SelectionField) CodegenMergeExpression() (string, error) {
	name := f.CodegenName()

	if f.IsAggregate() {
		switch f.Aggregation {
		case "sum":
			return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", name, name, name), nil
		case "min":
			return fmt.Sprintf("if(result.%s == nil or result.%s > data.%s) then result.%s = data.%s end", name, name, name, name, name), nil
		case "max":
			return fmt.Sprintf("if(result.%s == nil or result.%s < data.%s) then result.%s = data.%s end", name, name, name, name, name), nil
		case "count":
			if f.Distinct {
				return fmt.Sprintf("result.%s = sky_distinct_merge(result.%s, data.%s)", name, name, name), nil
			} else {
				return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", name, name, name), nil
			}
		case "avg":
			return fmt.Sprintf("result.%s = sky_average_merge(result.%s, data.%s)", name, name, name), nil
		case "histogram":
			return fmt.Sprintf("result.%s = sky_histogram_merge(result.%s, data.%s)", name, name, name), nil
		case "":
			return fmt.Sprintf("result.%s = data.%s", name, name), nil
		}
		return "", fmt.Errorf("SelectionField: Unsupported merge aggregation method: %s", f.Aggregation)
	}

	// Non-aggregate selection fields are handled by the parent selection.
	return "", nil
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
// Finalization
//--------------------------------------

// Finalizes the results into a final state after merge.
func (f *SelectionField) Finalize(data interface{}) error {
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
