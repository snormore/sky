package ast

import (
	"bytes"
	"fmt"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/db"
	"sort"
	"strings"
)

// A Query is a structured way of aggregating data in the database.
type Query struct {
	table             *core.Table
	factorizer        db.Factorizer
	source            string
	refs              []*VarRef
	sequence          int
	Prefix            string
	systemVariables   []*Variable
	declaredVariables []*Variable
	statements        Statements
}

// NewQuery returns a new query.
func NewQuery() *Query {
	q := &Query{statements: make(Statements, 0)}
	q.addSystemVariable(NewVariable("@eos", core.BooleanDataType))
	q.addSystemVariable(NewVariable("@eof", core.BooleanDataType))
	return q
}

// Retrieves the table this query is associated with.
func (q *Query) Table() *core.Table {
	return q.table
}

func (q *Query) SetTable(table *core.Table) {
	q.table = table
}

// Retrieves the factors database this query is associated with.
func (q *Query) Factorizer() db.Factorizer {
	return q.factorizer
}

func (q *Query) SetFactorizer(factorizer db.Factorizer) {
	q.factorizer = factorizer
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

// Adds a system variable.
func (q *Query) addSystemVariable(variable *Variable) {
	variable.SetParent(q)
	q.systemVariables = append(q.systemVariables, variable)
}

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
	// Check against system variables first.
	for _, v := range q.systemVariables {
		if v.Name == name {
			return v
		}
	}

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

	// Find all variables declared within statements.
	for _, v := range q.statements.Variables() {
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
// Code Generation
//--------------------------------------

// Generates Lua code for the query. Code can only generated once for the
// query and then it is cached.
func (q *Query) Codegen() (string, error) {
	if len(q.source) == 0 {
		buffer := new(bytes.Buffer)

		// Generate initialization functions (if necessary).
		if q.RequiresInitialization() {
			buffer.WriteString("----------------------------------------------------------------------\n")
			buffer.WriteString("--\n")
			buffer.WriteString("-- INITIALIZATION\n")
			buffer.WriteString("--\n")
			buffer.WriteString("----------------------------------------------------------------------\n\n")

			str, err := q.statements.CodegenAggregateFunctions(true)
			if err != nil {
				return "", err
			}
			buffer.WriteString(str)
			buffer.WriteString(q.CodegenAggregateFunction(true))
		}

		// Generate aggregation functions.
		buffer.WriteString("----------------------------------------------------------------------\n")
		buffer.WriteString("--\n")
		buffer.WriteString("-- AGGREGATION\n")
		buffer.WriteString("--\n")
		buffer.WriteString("----------------------------------------------------------------------\n\n")
		str, err := q.statements.CodegenAggregateFunctions(false)
		if err != nil {
			return "", err
		}
		buffer.WriteString(str)
		buffer.WriteString(q.CodegenAggregateFunction(false))

		// Generate merge functions.
		buffer.WriteString("----------------------------------------------------------------------\n")
		buffer.WriteString("--\n")
		buffer.WriteString("-- MERGE\n")
		buffer.WriteString("--\n")
		buffer.WriteString("----------------------------------------------------------------------\n\n")
		str, err = q.statements.CodegenMergeFunctions(map[string]interface{}{})
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

	// Begin cursor loop.
	fmt.Fprintln(buffer, "  repeat")

	// Call each statement function.
	for _, statement := range q.statements {
		fmt.Fprintf(buffer, "      %s(cursor, data)\n", statement.FunctionName(init))
	}

	// End cursor loop.
	fmt.Fprintln(buffer, "  until not cursor:next()")

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

// Converts factorized results from the aggregate function results to use
// the appropriate strings.
func (q *Query) Defactorize(data interface{}) error {
	return q.statements.Defactorize(data)
}

// Finalizes the results into a final state after merge.
func (q *Query) Finalize(data interface{}) error {
	return q.statements.Finalize(data)
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
