package ast

import (
	"fmt"
)

// VarRef represents a reference to a variable in the current scope in a query.
type VarRef struct {
	Value string
}

func (v *VarRef) node() string {}

// NewVarRef creates a new VarRef instance.
func NewVarRef() *VarRef {
	return &VarRef{}
}

func (v *VarRef) String() string {
	return "@" + v.value
}
