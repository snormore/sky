package ast

import (
	"bytes"
	"errors"
	"fmt"
)

// An Exit statement stops the query on the current object.
type Exit struct {
	queryElementImpl
}

// Creates a new exit.
func NewExit() *Exit {
	return &Exit{}
}

// Retrieves the function name used during codegen.
func (e *Exit) FunctionName(init bool) string {
	return fmt.Sprintf("a%d", e.ElementId())
}

// Retrieves the merge function name used during codegen.
func (e *Exit) MergeFunctionName() string {
	return ""
}

// Returns a list of variable references within this statement.
func (e *Exit) VarRefs() []*VarRef {
	return []*VarRef{}
}

// Returns a list of variables declared within this statement.
func (e *Exit) Variables() []*Variable {
	return []*Variable{}
}

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query selection into an untyped map.
func (e *Exit) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"type": TypeExit,
	}
}

// Decodes an exit from an untyped map.
func (e *Exit) Deserialize(obj map[string]interface{}) error {
	if obj == nil {
		return errors.New("Exit: Unable to deserialize nil.")
	}
	if obj["type"] != TypeExit {
		return fmt.Errorf("Exit: Invalid statement type: %v", obj["type"])
	}

	return nil
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the selection aggregation.
func (e *Exit) CodegenAggregateFunction(init bool) (string, error) {
	buffer := new(bytes.Buffer)

	fmt.Fprintf(buffer, "-- %s\n", e.String())
	fmt.Fprintf(buffer, "function %s(cursor, data)\n", e.FunctionName(init))
	fmt.Fprintf(buffer, "  error(exit_error)\n")
	fmt.Fprintln(buffer, "end")

	return buffer.String(), nil
}

// Generates Lua code for the selection merge.
func (e *Exit) CodegenMergeFunction(fields map[string]interface{}) (string, error) {
	return "", nil
}

func (e *Exit) Defactorize(data interface{}) error {
	return nil
}

func (e *Exit) Finalize(data interface{}) error {
	return nil
}

func (e *Exit) RequiresInitialization() bool {
	return false
}

// Converts the statement to a string-based representation.
func (e *Exit) String() string {
	return "EXIT"
}
