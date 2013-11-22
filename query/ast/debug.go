package ast

import (
	"bytes"
	"fmt"
)

// A Debug statement prints an expression to STDOUT.
type Debug struct {
	queryElementImpl
	expression Expression
}

// Creates a new debug statement.
func NewDebug() *Debug {
	return &Debug{}
}

// Retrieves the function name used during codegen.
func (d *Debug) FunctionName(init bool) string {
	return fmt.Sprintf("a%d", d.ElementId())
}

// Retrieves the merge function name used during codegen.
func (d *Debug) MergeFunctionName() string {
	return ""
}

// Returns the expression to be printed.
func (d *Debug) Expression() Expression {
	return d.expression
}

// Sets the expression to be printed.
func (d *Debug) SetExpression(expression Expression) {
	if d.expression != nil {
		d.expression.SetParent(nil)
	}
	d.expression = expression
	if d.expression != nil {
		d.expression.SetParent(d)
	}
}

// Returns a list of variable references within this debug statement.
func (d *Debug) VarRefs() []*VarRef {
	refs := []*VarRef{}
	refs = append(refs, d.expression.VarRefs()...)
	return refs
}

// Returns a list of variables declared within this statement.
func (d *Debug) Variables() []*Variable {
	return []*Variable{}
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the selection aggregation.
func (d *Debug) CodegenAggregateFunction(init bool) (string, error) {
	buffer := new(bytes.Buffer)

	expressionCode, err := d.expression.Codegen()
	if err != nil {
		return "", err
	}

	fmt.Fprintf(buffer, "-- %s\n", d.String())
	fmt.Fprintf(buffer, "function %s(cursor, data)\n", d.FunctionName(init))
	fmt.Fprintf(buffer, "  print(%s, cursor.next_event:eof())\n", expressionCode)
	fmt.Fprintln(buffer, "end")

	return buffer.String(), nil
}

// Generates Lua code for the selection merge.
func (d *Debug) CodegenMergeFunction(fields map[string]interface{}) (string, error) {
	return "", nil
}

func (d *Debug) Defactorize(data interface{}) error {
	return nil
}

func (d *Debug) Finalize(data interface{}) error {
	return nil
}

func (d *Debug) RequiresInitialization() bool {
	return false
}

// Converts the statement to a string-based representation.
func (d *Debug) String() string {
	return fmt.Sprintf("DEBUG(%s)", d.expression.String())
}
