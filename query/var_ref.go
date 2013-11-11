package query

import (
	"fmt"
)

type VarRef struct {
	queryElementImpl
	value string
}

// Retrieves the referenced variable.
func (v *VarRef) Variable() (*Variable, error) {
	q := v.Query()
	if q == nil {
		return nil, fmt.Errorf("Var ref not attached to query: %s", v.value)
	}
	variable := q.GetVariable(v.value)
	if variable == nil {
		return nil, fmt.Errorf("Variable not found: %s", v.value)
	}
	return variable, nil
}

func (v *VarRef) VarRefs() []*VarRef {
	return []*VarRef{v}
}

// Generates a Lua representation of this variable reference.
func (v *VarRef) Codegen() (string, error) {
	if _, err := v.Variable(); err != nil {
		return "", err
	}

	// Remove ampersand prefix if this is a system variable.
	if v.value == "@eof" {
		return "cursor:next_event:eof()", nil
	} else if v.value == "@eos" {
		return "cursor:event:eos()", nil
	} else {
		return fmt.Sprintf("cursor.event:%s()", v.value), nil
	}
}

// Generates a Lua representation of the raw struct reference.
func (v *VarRef) CodegenRaw() (string, error) {
	if _, err := v.Variable(); err != nil {
		return "", err
	}
	return fmt.Sprintf("cursor.event._%s", v.value), nil
}

func (v *VarRef) String() string {
	return "@" + v.value
}
