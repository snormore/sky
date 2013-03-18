package skyd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Query is a structured way of aggregating data in the database.
type Query struct {
	table           *Table
	sequence        int
	Steps           QueryStepList
	SessionIdleTime int
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewQuery returns a new query.
func NewQuery(table *Table) *Query {
	return &Query{
		table: table,
		Steps: make(QueryStepList, 0),
	}
}

//------------------------------------------------------------------------------
//
// Accessors
//
//------------------------------------------------------------------------------

// Retrieves the table this query is associated with.
func (q *Query) Table() *Table {
	return q.table
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query into an untyped map.
func (q *Query) Serialize() map[string]interface{} {
	obj := map[string]interface{}{
		"sessionIdleTime": q.SessionIdleTime,
		"steps":           q.Steps.Serialize(),
	}
	return obj
}

// Decodes a query from an untyped map.
func (q *Query) Deserialize(obj map[string]interface{}) error {
	var err error

	// Deserialize "session idle time".
	if sessionIdleTime, ok := obj["sessionIdleTime"].(float64); ok || obj["sessionIdleTime"] == nil {
		q.SessionIdleTime = int(sessionIdleTime)
	} else {
		return fmt.Errorf("Invalid 'sessionIdleTime': %v", obj["sessionIdleTime"])
	}

	q.Steps, err = DeserializeQueryStepList(obj["steps"], q)
	if err != nil {
		return err
	}
	return nil
}

//--------------------------------------
// Encoding
//--------------------------------------

// Encodes a query to JSON.
func (q *Query) Encode(writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(q.Serialize())
	return err
}

// Decodes a query from JSON.
func (q *Query) Decode(reader io.Reader) error {
	// Decode into an untyped object first since we need to determine the
	// type of steps to create.
	var obj map[string]interface{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&obj)
	if err != nil {
		return err
	}

	return q.Deserialize(obj)
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates Lua code for the query.
func (q *Query) Codegen() (string, error) {
	buffer := new(bytes.Buffer)
	err := q.codegenSteps(buffer)
	if err != nil {
		return "", err
	}
	q.codegenAggregateFunction(buffer)
	q.codegenMergeFunction(buffer)
	return buffer.String(), nil
}

// Recursively generates code for all steps.
func (q *Query) codegenSteps(writer io.Writer) error {
	for _, step := range q.Steps {
		code, err := step.Codegen()
		if err != nil {
			return err
		}
		fmt.Fprintln(writer, code)
	}
	return nil
}

// Generates the 'aggregate()' function.
func (q *Query) codegenAggregateFunction(writer io.Writer) {
	// Generate the function definition.
	fmt.Fprintln(writer, "function aggregate(cursor, data)")
	
	// Set the session idle if one is available.
	if q.SessionIdleTime > 0 {
		fmt.Fprintf(writer, "  cursor:set_session_idle(%d)\n", q.SessionIdleTime)
	}
	
	// Begin cursor loop.
	fmt.Fprintln(writer, "  while cursor:next_session() do")
	fmt.Fprintln(writer, "    while cursor:next() do")

	// Call each step function.
	for _, step := range q.Steps {
		fmt.Fprintf(writer, "      %s(cursor, data)\n", step.FunctionName())
	}
	
	// End cursor loop.
	fmt.Fprintln(writer, "    end")
	fmt.Fprintln(writer, "  end")

	// End function.
	fmt.Fprintln(writer, "end\n")
}

// Generates the 'merge()' function.
func (q *Query) codegenMergeFunction(writer io.Writer) {
	// Generate the function definition.
	fmt.Fprintln(writer, "function merge(results, data)")
	
	// Call each step function if it has a merge function.
	for _, step := range q.Steps {
		if step.MergeFunctionName() != "" {
			fmt.Fprintf(writer, "  %s(results, data)\n", step.MergeFunctionName())
		}
	}
	
	// End function.
	fmt.Fprintln(writer, "end\n")
}

// Returns an autoincrementing numeric identifier.
func (q *Query) NextIdentifier() int {
	q.sequence += 1
	return q.sequence
}
