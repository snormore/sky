package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/factors"
	"io"
)

// A Query is a structured way of aggregating data in the database.
type Query struct {
	table           *core.Table
	fdb             *factors.DB
	sequence        int
	Prefix          string
	statements      Statements
	SessionIdleTime int
}

// NewQuery returns a new query.
func NewQuery() *Query {
	return &Query{statements: make(Statements, 0)}
}

// Retrieves the table this query is associated with.
func (q *Query) Table() *core.Table {
	return q.table
}

func (q *Query) SetTable(table *core.Table) {
	q.table = table
}

// Retrieves the factors database this query is associated with.
func (q *Query) Fdb() *factors.DB {
	return q.fdb
}

func (q *Query) SetFdb(fdb *factors.DB) {
	q.fdb = fdb
}

// Returns the top-level statements of the query.
func (q *Query) Statements() Statements {
	return q.statements
}

// Sets the condition's statements.
func (q *Query) SetStatements(statements Statements) {
	for _, s := range q.statements {
		s.SetParent(nil)
	}
	q.statements = statements
	for _, s := range q.statements {
		s.SetParent(q)
	}
}

//--------------------------------------
// QueryElement interface
//--------------------------------------

func (q *Query) Parent() QueryElement {
	return nil
}

func (q *Query) SetParent(e QueryElement) {
}

func (q *Query) Query() *Query {
	return q
}

func (q *Query) ElementId() int {
	return 0
}

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a query into an untyped map.
func (q *Query) Serialize() map[string]interface{} {
	obj := map[string]interface{}{
		"sessionIdleTime": q.SessionIdleTime,
		"prefix":          q.Prefix,
		"statements":      q.statements.Serialize(),
	}
	return obj
}

// Decodes a query from an untyped map.
func (q *Query) Deserialize(obj map[string]interface{}) error {
	// Deserialize "prefix".
	if prefix, ok := obj["prefix"].(string); ok || obj["prefix"] == nil {
		q.Prefix = prefix
	} else {
		return fmt.Errorf("Invalid 'prefix': %v", obj["prefix"])
	}

	// Deserialize "session idle time".
	if sessionIdleTime, ok := obj["sessionIdleTime"].(float64); ok || obj["sessionIdleTime"] == nil {
		q.SessionIdleTime = int(sessionIdleTime)
	} else {
		return fmt.Errorf("Invalid 'sessionIdleTime': %v", obj["sessionIdleTime"])
	}

	// Parse the string-based query if it's passed in.
	q.SetStatements(Statements{})
	if content, ok := obj["content"].(string); ok && len(content) > 0 {
		tmp, err := NewParser().ParseString(content)
		if err != nil {
			return err
		}
		q.SetStatements(tmp.Statements())
	}

	if len(q.statements) == 0 {
		statements, err := DeserializeStatements(obj["statements"])
		if err != nil {
			return err
		}
		q.SetStatements(statements)
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
	// type of statements to create.
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

	// Generate initialization functions (if necessary).
	if q.RequiresInitialization() {
		str, err := q.statements.CodegenAggregateFunctions(true)
		if err != nil {
			return "", err
		}
		buffer.WriteString(str)
		buffer.WriteString(q.CodegenAggregateFunction(true))
	}

	// Generate aggregation functions.
	str, err := q.statements.CodegenAggregateFunctions(false)
	if err != nil {
		return "", err
	}
	buffer.WriteString(str)
	buffer.WriteString(q.CodegenAggregateFunction(false))

	// Generate merge functions.
	str, err = q.statements.CodegenMergeFunctions()
	if err != nil {
		return "", err
	}
	buffer.WriteString(str)
	buffer.WriteString(q.CodegenMergeFunction())

	return buffer.String(), nil
}

// Generates the 'aggregate()' function. If the "init" flag is specified then
// a special set of initialization functions are created to generate the data
// structure used by the main aggregate function. This is primarily used to
// setup histogram structures.
func (q *Query) CodegenAggregateFunction(init bool) string {
	buffer := new(bytes.Buffer)

	// Generate the function definition.
	if init {
		fmt.Fprintln(buffer, "function initialize(cursor, data)")
	} else {
		fmt.Fprintln(buffer, "function aggregate(cursor, data)")
	}

	// Set the session idle if one is available.
	if q.SessionIdleTime > 0 {
		fmt.Fprintf(buffer, "  cursor:set_session_idle(%d)\n", q.SessionIdleTime)
	}

	// Begin cursor loop.
	fmt.Fprintln(buffer, "  while cursor:next_session() do")
	fmt.Fprintln(buffer, "    while cursor:next() do")

	// Call each statement function.
	for _, statement := range q.statements {
		fmt.Fprintf(buffer, "      %s(cursor, data)\n", statement.FunctionName(init))
	}

	// End cursor loop.
	fmt.Fprintln(buffer, "    end")
	fmt.Fprintln(buffer, "  end")

	// End function.
	fmt.Fprintln(buffer, "end\n")

	return buffer.String()
}

// Generates the 'merge()' function.
func (q *Query) CodegenMergeFunction() string {
	buffer := new(bytes.Buffer)

	// Generate the function definition.
	fmt.Fprintln(buffer, "function merge(results, data)")

	// Call each statement function if it has a merge function.
	fmt.Fprintf(buffer, q.statements.CodegenMergeInvoke())

	// End function.
	fmt.Fprintln(buffer, "end\n")

	return buffer.String()
}

// Returns an autoincrementing numeric identifier.
func (q *Query) NextIdentifier() int {
	q.sequence += 1
	return q.sequence
}

//--------------------------------------
// Factorization
//--------------------------------------

// Converts factorized results from the aggregate function results to use
// the appropriate strings.
func (q *Query) Defactorize(data interface{}) error {
	return q.statements.Defactorize(data)
}

//--------------------------------------
// Initialization
//--------------------------------------

// Checks if this query requires a data structure to be initialized before
// performing aggregation. This function returns true if any nested query
// statements require initialization.
func (q *Query) RequiresInitialization() bool {
	return q.statements.RequiresInitialization()
}

//--------------------------------------
// String
//--------------------------------------

// Convert the query to a string-based representation.
func (q *Query) String() string {
	var str string
	str += q.statements.String()
	return str
}
