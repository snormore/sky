package ast

import (
	"fmt"
	"github.com/skydb/sky/core"
)

// VarDecl represents a variable declaration in the query. The value
// of the variable persist for the duration of an object and can be
// referenced like any other property on the database. The variable can
// also be associated with another variable for the purpose of reusing
// factorization.
type VarDecl struct {
	Id          int64
	Name        string
	DataType    string
	Association string
}

func (v *VarDecl) node() {}

// NewVarDecl returns a new VarDecl instance.
func NewVarDecl(name string, dataType string) *VarDecl {
	return &VarDecl{Name: name, DataType: dataType}
}

// Determines if the variable declaration is a system variable.
func (v *VarDecl) IsSystemVarDecl() bool {
	return (len(v.Name) != 0 && v.Name[0] == '@')
}

func (v *VarDecl) String() string {
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
