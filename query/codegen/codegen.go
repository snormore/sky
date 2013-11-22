package codegen

import (
	"bytes"
	"io"

	"github.com/skydb/sky/query/ast"
)

// Codegen generates the Lua code required to execute a given query.
func Codegen(w io.Writer, q *ast.Query) error {
	return tmpl.Execute(w, q)
}

// CodegenString generates a query's Lua code to a string.
func CodegenString(q *ast.Query) (string, error) {
	var b bytes.Buffer
	if err := Codegen(&b, q); err != nil {
		return "", err
	}
	return b.String(), nil
}
