package query

import (
	"fmt"
	"github.com/skydb/sky/core"
)

type VarRef struct {
	queryElementImpl
	value string
}

// Retrieves the property that the variable references.
func (v *VarRef) Property() (*core.Property, error) {
	property := v.Query().table.PropertyFile().GetPropertyByName(v.value)
	if property == nil {
		return nil, fmt.Errorf("Property not found: %s", v.value)
	}
	return property, nil
}

// Generates a Lua representation of this variable reference.
func (v *VarRef) Codegen() (string, error) {
	_, err := v.Property()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("cursor.event:%s()", v.value), nil
}

func (v *VarRef) String() string {
	return v.value
}
