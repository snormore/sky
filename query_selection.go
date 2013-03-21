package skyd

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A selection step aggregates data in a query.
type QuerySelection struct {
	query             *Query
	functionName      string
	mergeFunctionName string
	Expression        string
	Alias             string
	Dimensions        []string
	Steps             QueryStepList
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
		query:             query,
		functionName:      fmt.Sprintf("a%d", id),
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

// Retrieves the child steps.
func (s *QuerySelection) GetSteps() QueryStepList {
	return s.Steps
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

// Generates Lua code for the selection aggregation.
func (s *QuerySelection) CodegenAggregateFunction() (string, error) {
	buffer := new(bytes.Buffer)

	// Generate child steps.
	str, err := s.Steps.CodegenAggregateFunctions()
	if err != nil {
		return "", err
	}
	buffer.WriteString(str + "\n")

	// Generate main function.
	fmt.Fprintf(buffer, "function %s(cursor, data)\n", s.FunctionName())

	// Group by dimension.
	for _, dimension := range s.Dimensions {
		fmt.Fprintf(buffer, "  dimension = cursor.event:%s()\n", dimension)
		fmt.Fprintf(buffer, "  if data.%s == nil then data.%s = {} end\n", dimension, dimension)
		fmt.Fprintf(buffer, "  if data.%s[dimension] == nil then data.%s[dimension] = {} end\n", dimension, dimension)
		fmt.Fprintf(buffer, "  data = data.%s[dimension]\n\n", dimension)
	}

	// Select value.
	exp, err := s.CodegenExpression()
	if err != nil {
		return "", err
	}
	fmt.Fprintln(buffer, "  "+exp)

	// End function definition.
	fmt.Fprintln(buffer, "end")

	return buffer.String(), nil
}

// Generates Lua code for the selection merge.
func (s *QuerySelection) CodegenMergeFunction() (string, error) {
	buffer := new(bytes.Buffer)

	// Generate child steps.
	str, err := s.Steps.CodegenMergeFunctions()
	if err != nil {
		return "", err
	}
	buffer.WriteString(str + "\n")

	// Generate nested functions first.
	code, err := s.CodegenInnerMergeFunction(0)
	if err != nil {
		return "", err
	}
	buffer.WriteString(code + "\n")

	// Generate main function.
	fmt.Fprintf(buffer, "function %s(result, data)\n", s.MergeFunctionName())
	fmt.Fprintf(buffer, "  %sn0(result, data)\n", s.MergeFunctionName())
	fmt.Fprintf(buffer, "end\n")

	return buffer.String(), nil
}

// Generates Lua code for the inner merge.
func (s *QuerySelection) CodegenInnerMergeFunction(index int) (string, error) {
	buffer := new(bytes.Buffer)

	// Generate next nested function first.
	if index < len(s.Dimensions) {
		code, err := s.CodegenInnerMergeFunction(index + 1)
		if err != nil {
			return "", err
		}
		buffer.WriteString(code + "\n")
	}

	// Generate a rollup if our index points at a dimension. Otherwise generate
	// the leaf merge.
	fmt.Fprintf(buffer, "function %sn%d(result, data)\n", s.MergeFunctionName(), index)
	if index < len(s.Dimensions) {
		dimension := s.Dimensions[index]
		fmt.Fprintf(buffer, "  if data ~= nil and data.%s ~= nil then\n", dimension)
		fmt.Fprintf(buffer, "    if result.%s == nil then result.%s = {} end\n", dimension, dimension)
		fmt.Fprintf(buffer, "    for k,v in pairs(data.%s) do\n", dimension)
		fmt.Fprintf(buffer, "      if result.%s[k] == nil then result.%s[k] = {} end\n", dimension, dimension)
		fmt.Fprintf(buffer, "      %sn%d(result.%s[k], v)\n", s.MergeFunctionName(), (index + 1), dimension)
		fmt.Fprintf(buffer, "    end\n")
		fmt.Fprintf(buffer, "  end\n")
	} else {
		// Merge value.
		exp, err := s.CodegenMergeExpression()
		if err != nil {
			return "", err
		}
		fmt.Fprintln(buffer, "  "+exp)
	}
	fmt.Fprintf(buffer, "end\n")

	return buffer.String(), nil
}

// Generates Lua code for the expression.
func (s *QuerySelection) CodegenExpression() (string, error) {
	r, _ := regexp.Compile(`^ *(?:count\(\)|(sum|min|max)\((\w+)\)|(\w+)) *$`)
	if m := r.FindStringSubmatch(s.Expression); m != nil {
		if len(m[1]) > 0 { // sum()/min()/max()
			switch m[1] {
			case "sum":
				return fmt.Sprintf("data.%s = (data.%s or 0) + cursor.event:%s()", s.Alias, s.Alias, m[2]), nil
			case "min":
				return fmt.Sprintf("if(data.%s == nil or data.%s > cursor.event:%s()) then data.%s = cursor.event:%s() end", s.Alias, s.Alias, m[2], s.Alias, m[2]), nil
			case "max":
				return fmt.Sprintf("if(data.%s == nil or data.%s < cursor.event:%s()) then data.%s = cursor.event:%s() end", s.Alias, s.Alias, m[2], s.Alias, m[2]), nil
			}
		} else if len(m[3]) > 0 { // assignment
			return fmt.Sprintf("data.%s = cursor.event:%s()", s.Alias, m[3]), nil
		} else { // count()
			return fmt.Sprintf("data.%s = (data.%s or 0) + 1", s.Alias, s.Alias), nil
		}
	}

	return "", fmt.Errorf("skyd.QuerySelection: Invalid expression: %q", s.Expression)
}

// Generates Lua code for the merge expression.
func (s *QuerySelection) CodegenMergeExpression() (string, error) {
	r, _ := regexp.Compile(`^ *(?:count\(\)|(sum|min|max)\((\w+)\)|(\w+)) *$`)
	if m := r.FindStringSubmatch(s.Expression); m != nil {
		if len(m[1]) > 0 { // sum()/min()/max()
			switch m[1] {
			case "sum":
				return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", s.Alias, s.Alias, s.Alias), nil
			case "min":
				return fmt.Sprintf("if(result.%s == nil or result.%s > data.%s) then result.%s = data.%s end", s.Alias, s.Alias, s.Alias, s.Alias, s.Alias), nil
			case "max":
				return fmt.Sprintf("if(result.%s == nil or result.%s < data.%s) then result.%s = data.%s end", s.Alias, s.Alias, s.Alias, s.Alias, s.Alias), nil
			}
		} else if len(m[3]) > 0 { // assignment
			return fmt.Sprintf("result.%s = data.%s", s.Alias, s.Alias), nil
		} else { // count()
			return fmt.Sprintf("result.%s = (result.%s or 0) + (data.%s or 0)", s.Alias, s.Alias, s.Alias), nil
		}
	}

	return "", fmt.Errorf("skyd.QuerySelection: Invalid merge expression: %q", s.Expression)
}

//--------------------------------------
// Factorization
//--------------------------------------

// Converts factorized fields back to their original strings.
func (s *QuerySelection) Defactorize(data interface{}) (error) {
	err := s.defactorize(data, 0)
	if err != nil {
		return err
	}

	// Defactorize child steps.
	err = s.Steps.Defactorize(data)

	return err
}

// Recursively defactorizes dimensions.
func (s *QuerySelection) defactorize(data interface{}, index int) (error) {
	if index >= len(s.Dimensions) {
		return nil
	}
	// Ignore any values that are nil or not maps.
	inner, ok := data.(map[interface{}]interface{})
	if !ok || data == nil {
		return nil
	}
	
	// Retrieve property.
	dimension := s.Dimensions[index]
	property := s.query.table.propertyFile.GetPropertyByName(dimension)
	if property == nil {
		return fmt.Errorf("skyd.QuerySelection: Property not found: %s", dimension)
	}
	
	// Defactorize.
	if outer, ok := inner[dimension].(map[interface{}]interface{}); ok {
		copy := map[interface{}]interface{}{}
		for k, v := range outer {
			if property.DataType == FactorDataType {
				sequence, err := castUint64(k)
				if err != nil {
					return err
				}
				stringValue, err := s.query.factors.Defactorize(s.query.table.Name(), dimension, sequence)
				if err != nil {
					return err
				}
				copy[stringValue] = v
			} else {
				copy[k] = v
			}

			// Defactorize next dimension.
			s.defactorize(v, index+1)
		}
		inner[dimension] = copy
	}
	
	return nil
}
