package query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/factors"
	"io"
	"sort"
	"strings"
)

// A Query is a structured way of aggregating data in the database.
type Query struct {
	table             *core.Table
	fdb               *factors.DB
	source            string
	refs              []*VarRef
	sequence          int
	Prefix            string
	SessionIdleTime   int
	declaredVariables []*Variable
	statements        Statements
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

// Sets the query's statements.
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
// Variables
//--------------------------------------

// Returns the variables declared for this query.
func (q *Query) DeclaredVariables() []*Variable {
	return q.declaredVariables
}

// Sets the variables declared on this query.
func (q *Query) SetDeclaredVariables(variables []*Variable) {
	for _, v := range q.declaredVariables {
		v.SetParent(nil)
	}
	q.declaredVariables = variables
	for _, v := range q.declaredVariables {
		v.SetParent(q)
	}
}

// Retrieves the variable with a given name.
func (q *Query) GetVariable(name string) *Variable {
	// Try to find variable from schema first.
	if q.Table() != nil && q.Table().PropertyFile() != nil {
		p := q.Table().PropertyFile().GetPropertyByName(name)
		if p != nil {
			return &Variable{
				Name:       p.Name,
				DataType:   p.DataType,
				PropertyId: p.Id,
			}
		}
	}

	// Otherwise check if the variable was declared.
	for _, v := range q.declaredVariables {
		if v.Name == name {
			return v
		}
	}

	return nil
}

// Returns a list of all variables explicitly declared or implicitly
// declared through the table schema.
func (q *Query) Variables() []*Variable {
	variables := []*Variable{}
	lookup := map[string]bool{}

	for _, ref := range q.VarRefs() {
		if !lookup[ref.value] {
			if v := q.GetVariable(ref.value); v != nil {
				variables = append(variables, v)
				lookup[ref.value] = true
			}
		}
	}

	sort.Sort(Variables(variables))

	return variables
}

// Retrieves a list of variables references by the query. This list is
// cached so it is only generated the first time it's called.
func (q *Query) VarRefs() []*VarRef {
	if q.refs == nil {
		q.refs = []*VarRef{}
		for _, s := range q.statements {
			q.refs = append(q.refs, s.VarRefs()...)
		}
	}
	return q.refs
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
	var err error

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

	// DEPRECATED: Statements can be passed in as "steps".
	val := obj["steps"]
	if val == nil {
		val = obj["statements"]
	}

	// Parse statements as string or as map.
	var statements Statements
	if strval, ok := val.(string); ok && len(strval) > 0 {
		if statements, err = NewStatementsParser().ParseString(strval); err != nil {
			return err
		}
	} else {
		if statements, err = DeserializeStatements(val); err != nil {
			return err
		}
	}
	q.SetStatements(statements)

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

// Generates Lua code for the query. Code can only generated once for the
// query and then it is cached.
func (q *Query) Codegen() (string, error) {
	if len(q.source) == 0 {
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

		q.source = buffer.String()
	}

	return q.source, nil
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
// Properties
//--------------------------------------

// Determines the minimum and maximum property identifiers that are used
// in this query.
func (q *Query) PropertyIdentifierRange() (int64, int64) {
	maxPropertyId, minPropertyId := q.Table().PropertyFile().NextIdentifiers()
	return minPropertyId, maxPropertyId
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
	arr := []string{}
	for _, v := range q.declaredVariables {
		arr = append(arr, v.String())
	}
	arr = append(arr, q.statements.String())
	return strings.Join(arr, "\n")
}
