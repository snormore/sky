package ast

import (
	"bytes"
	"fmt"
	"strings"
)

type Statements []Statement

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
func (s Statements) CodegenMergeFunctions(fields map[string]interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	for _, statement := range s {
		code, err := statement.CodegenMergeFunction(fields)
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
		if block_statement, ok := statement.(BlockStatement); ok {
			code := block_statement.Statements().CodegenMergeInvoke()
			if code != "" {
				fmt.Fprintf(buffer, code)
			}
		}
	}
	return buffer.String()
}

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

// Finalizes results into an end state.
func (s Statements) Finalize(data interface{}) error {
	for _, statement := range s {
		if err := statement.Finalize(data); err != nil {
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
// Utility
//--------------------------------------

// Retrieves a list of variable references used by the statements.
func (s Statements) VarRefs() []*VarRef {
	refs := []*VarRef{}
	for _, statement := range s {
		refs = append(refs, statement.VarRefs()...)
	}
	return refs
}

// Retrieves a list of variable declarations within the statements.
func (s Statements) Variables() []*Variable {
	variables := []*Variable{}
	for _, statement := range s {
		variables = append(variables, statement.Variables()...)
	}
	return variables
}

// Converts the statements to a string-based representation.
func (s Statements) String() string {
	output := []string{}
	for _, statement := range s {
		output = append(output, statement.String())
	}
	return strings.Join(output, "\n")
}
