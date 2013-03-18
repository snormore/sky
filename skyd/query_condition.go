package skyd

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	QueryConditionUnitSteps    = "steps"
	QueryConditionUnitSessions = "sessions"
	QueryConditionUnitSeconds  = "seconds"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A condition step made within a query.
type QueryCondition struct {
	query       *Query
	functionName string
	Expression  string
	Within      int
	WithinUnits string
	Steps       QueryStepList
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// Creates a new condition.
func NewQueryCondition(query *Query) *QueryCondition {
	id := query.NextIdentifier()
	return &QueryCondition{
		query:       query,
		functionName: fmt.Sprintf("f%d", id),
		Within:      0,
		WithinUnits: QueryConditionUnitSteps,
	}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// Retrieves the query this condition is associated with.
func (c *QueryCondition) Query() *Query {
	return c.query
}

// Retrieves the function name used during codegen.
func (c *QueryCondition) FunctionName() string {
	return c.functionName
}

// Retrieves the merge function name used during codegen.
func (c *QueryCondition) MergeFunctionName() string {
	return ""
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query condition into an untyped map.
func (c *QueryCondition) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type":        QueryStepTypeCondition,
		"expression":  c.Expression,
		"within":      c.Within,
		"withinUnits": c.WithinUnits,
		"steps":       c.Steps.Serialize(),
	}
}

// Decodes a query condition from an untyped map.
func (c *QueryCondition) Deserialize(obj map[string]interface{}) error {
	if obj == nil {
		return errors.New("skyd.QueryCondition: Unable to deserialize nil.")
	}
	if obj["type"] != QueryStepTypeCondition {
		return fmt.Errorf("skyd.QueryCondition: Invalid step type: %v", obj["type"])
	}

	// Deserialize "expression".
	if expression, ok := obj["expression"].(string); ok {
		c.Expression = expression
	} else {
		return fmt.Errorf("Invalid 'expression': %v", obj["expression"])
	}

	// Deserialize "within".
	if within, ok := obj["within"].(float64); ok {
		c.Within = int(within)
	} else {
		return fmt.Errorf("Invalid 'within': %v", obj["within"])
	}

	// Deserialize "within units".
	if withinUnits, ok := obj["withinUnits"].(string); ok {
		switch withinUnits {
		case QueryConditionUnitSteps, QueryConditionUnitSessions, QueryConditionUnitSeconds:
			c.WithinUnits = withinUnits
		default:
			return fmt.Errorf("Invalid 'within units': %v", withinUnits)
		}
	} else {
		return fmt.Errorf("Invalid 'within units': %v", obj["within"])
	}

	// Deserialize steps.
	var err error
	c.Steps, err = DeserializeQueryStepList(obj["steps"], c.query)
	if err != nil {
		return err
	}

	return nil
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the query.
func (c *QueryCondition) Codegen() (string, error) {
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "function %s(cursor, data)\n", c.FunctionName())
	fmt.Fprintf(buffer, "  if cursor:eos() or cursor:eof() then return false end\n")
	if c.WithinUnits == QueryConditionUnitSteps {
		fmt.Fprintf(buffer, "  remaining = %d", c.Within)
	}
	fmt.Fprintf(buffer, "  repeat\n")
	if c.WithinUnits == QueryConditionUnitSteps {
		fmt.Fprintf(buffer, "    if remaining <= 0 then return false end\n")
	}
	fmt.Fprintf(buffer, "    if %s then\n", c.CodegenExpression())

	// Call each step function.
	fmt.Println(len(c.Steps))
	for _, step := range c.Steps {
		fmt.Fprintf(buffer, "        %s(cursor, data)\n", step.FunctionName())
	}

	fmt.Fprintf(buffer, "      return true\n")
	fmt.Fprintf(buffer, "    end\n")
	if c.WithinUnits == QueryConditionUnitSteps {
		fmt.Fprintf(buffer, "    remaining = remaining - 1\n")
	}
	fmt.Fprintf(buffer, "  until not cursor:next()\n")
	fmt.Fprintf(buffer, "  return false\n")

	// End function definition.
	fmt.Fprintln(buffer, "end")
	
	return buffer.String(), nil
}

// Generates Lua code for the expression.
func (c *QueryCondition) CodegenExpression() string {
	r, _ := regexp.Compile(`^ *(\w+) *(==) *("[^"]*"|'[^']*'|\d+|true|false) *$`)
	m := r.FindStringSubmatch(c.Expression)
	if m == nil {
		return "false"
	}
	return fmt.Sprintf("cursor.event:%s() %s %s", m[1], m[2], m[3])
}