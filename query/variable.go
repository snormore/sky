package query

import (
	"fmt"
	"github.com/skydb/sky/core"
)

// A Variable represents a variable declaration on the cursor. The value
// of the variable persist for the duration of an object and can be
// referenced like any other property on the database. The variable can
// also be associated with another variable for the purpose of reusing
// factorization.
type Variable struct {
	queryElementImpl
	Name        string
	DataType    string
	Association string
	PropertyId  int64
}

func NewVariable(name string, dataType string) *Variable {
	return &Variable{Name: name, DataType: dataType}
}

// Determines if the variable is a system variable.
func (v *Variable) IsSystemVariable() bool {
	return (len(v.Name) != 0 && v.Name[0] == '@')
}

func (v *Variable) Codegen() (string, error) {
	return "", nil
}

func (v *Variable) String() string {
	var dataType string
	switch v.DataType {
	case core.FactorDataType:
		dataType = "FACTOR"
	case core.StringDataType:
		dataType = "STRING"
	case core.IntegerDataType:
		dataType = "INTEGER"
	case core.FloatDataType:
		dataType = "FLOAT"
	case core.BooleanDataType:
		dataType = "BOOLEAN"
	}

	str := fmt.Sprintf("DECLARE @%s AS %s", v.Name, dataType)
	if v.Association != "" {
		str += "(@" + v.Association + ")"
	}
	return str
}
