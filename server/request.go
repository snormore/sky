package server

import (
	"github.com/skydb/sky/core"
)

// Request is a high level interface to the router methods on the server.
// It abstracts HTTP request to a more basic level.
type Request interface {
	Var(string) string
	Data() interface{}
	Table() *core.Table
	SetTable(*core.Table)
	Property() *core.Property
	SetProperty(*core.Property)
}

// request is the concrete implementation of the Request interface.
type request struct {
	vars     map[string]string
	data     interface{}
	table    *core.Table
	property *core.Property
}

// Var returns a URL string variable.
func (r *request) Var(key string) string {
	return r.vars[key]
}

// Data returns the parsed input data.
func (r *request) Data() interface{} {
	return r.data
}

// Table returns the table opened by the request.
func (r *request) Table() *core.Table {
	return r.table
}

// SetTable sets the table for the request.
func (r *request) SetTable(t *core.Table) {
	r.table = t
}

// Property returns the property associated with the request.
func (r *request) Property() *core.Property {
	return r.property
}

// SetProperty sets the property for the request.
func (r *request) SetProperty(t *core.Property) {
	r.property = t
}
