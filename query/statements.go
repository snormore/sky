package query

import (
	"bytes"
	"fmt"
	"strings"
)

type Statements []Statement

//--------------------------------------
// Serialization
//--------------------------------------

// Encodes a list of query statement into an untyped object.
func (s Statements) Serialize() []interface{} {
	statements := make([]interface{}, 0)
	for _, statement := range s {
		statements = append(statements, statement.Serialize())
	}
	return statements
}

// Decodes a list of query statements from an untyped slice.
func DeserializeStatements(obj interface{}, q *Query) (Statements, error) {
	l := make(Statements, 0)
	if statements, ok := obj.([]interface{}); ok {
		for _, _s := range statements {
			if s, ok := _s.(map[string]interface{}); ok {
				var statement Statement
				switch s["type"] {
				case TypeCondition:
					statement = NewCondition(q)
				case TypeSelection:
					statement = NewSelection(q)
				default:
					return nil, fmt.Errorf("Invalid query statement type: %v", s["type"])
				}
				err := statement.Deserialize(s)
				if err != nil {
					return nil, err
				}
				l = append(l, statement)
			} else {
				return nil, fmt.Errorf("Invalid statement: %v", obj)
			}
		}
	} else if obj != nil {
		return nil, fmt.Errorf("Invalid statements: %v", obj)
	}
	return l, nil
}

//--------------------------------------
// Code Generation
//--------------------------------------

// Generates aggregate code for all statements.
func (s Statements) CodegenAggregateFunctions(init bool) (string, error) {
	buffer := new(bytes.Buffer)
	for _, statement := range s {
		code, err := statement.CodegenAggregateFunction(init)
		if err != nil {
			return "", err
		}
		fmt.Fprintln(buffer, code)
	}
	return buffer.String(), nil
}

// Generates merge code for all statements.
func (s Statements) CodegenMergeFunctions() (string, error) {
	buffer := new(bytes.Buffer)
	for _, statement := range s {
		code, err := statement.CodegenMergeFunction()
		if err != nil {
			return "", err
		}
		fmt.Fprintln(buffer, code)
	}
	return buffer.String(), nil
}

// Generates merge invocations.
func (s Statements) CodegenMergeInvoke() string {
	buffer := new(bytes.Buffer)
	for _, statement := range s {
		// Generate this statement's invocation if available.
		if statement.MergeFunctionName() != "" {
			fmt.Fprintf(buffer, "  %s(results, data)\n", statement.MergeFunctionName())
		}

		// Recursively generate child statement invocations.
		code := statement.GetStatements().CodegenMergeInvoke()
		if code != "" {
			fmt.Fprintf(buffer, code)
		}
	}
	return buffer.String()
}

//--------------------------------------
// Factorization
//--------------------------------------

// Defactorizes results generated from the aggregate function.
func (s Statements) Defactorize(data interface{}) error {
	for _, statement := range s {
		err := statement.Defactorize(data)
		if err != nil {
			return err
		}
	}
	return nil
}

//--------------------------------------
// Initialization
//--------------------------------------

// Checks if this statement requires a data structure to be initialized before
// performing aggregation. This function returns true if any statements require
// initialization.
func (s Statements) RequiresInitialization() bool {
	for _, statement := range s {
		if statement.RequiresInitialization() {
			return true
		}
	}
	return false
}

//--------------------------------------
// String
//--------------------------------------

// Converts the statements to a string-based representation.
func (s Statements) String() string {
	output := []string{}
	for _, statement := range s {
		output = append(output, statement.String())
	}
	return strings.Join(output, "\n")
}
