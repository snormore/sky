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
	variable := v.Query().GetVariable(v.value)
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
	return fmt.Sprintf("cursor.event:%s()", v.value), nil
}

// Generates a Lua representation of the raw struct reference.
func (v *VarRef) RawCodegen() (string, error) {
	if _, err := v.Variable(); err != nil {
		return "", err
	}
	return fmt.Sprintf("cursor.event._%s", v.value), nil
}

func (v *VarRef) String() string {
	return v.value
}
