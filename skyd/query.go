package skyd

import (
	"fmt"
	"encoding/json"
	"io"
)

//------------------------------------------------------------------------------
//
// Constants
//
//------------------------------------------------------------------------------

const (
	QueryStepTypeCondition = "condition"
	QueryStepTypeSelection = "selection"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Query is a structured way of aggregating data in the database.
type Query struct {
	table *Table
	Steps []QueryStep
}

type QueryStep interface {
	Serialize() map[string]interface{}
	Deserialize(map[string]interface{}) error
	Codegen() (string, error)
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
		Steps: make([]QueryStep, 0),
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
	obj := map[string]interface{}{"steps": []interface{}{}}
	steps := make([]interface{}, 0)
	for _, step := range q.Steps {
		steps = append(steps, step.Serialize())
	}
	obj["steps"] = steps
	return obj
}

// Decodes a query from an untyped map.
func (q *Query) Deserialize(obj map[string]interface{}) error {
	q.Steps = make([]QueryStep, 0)
	if steps, ok := obj["steps"].([]map[string]interface{}); ok {
		for _, s := range steps  {
			var step QueryStep
			switch s["type"] {
			case QueryStepTypeCondition:
				step = NewQueryCondition(q)
			case QueryStepTypeSelection:
				step = NewQuerySelection(q)
			default:
				return fmt.Errorf("Invalid query step type: %v", s["type"])
			}
			step.Deserialize(s)
			q.Steps = append(q.Steps, step)
		}
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
	return "", nil
}
