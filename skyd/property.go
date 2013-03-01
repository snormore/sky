package skyd

import (
  "fmt"
)

// A Property is a loose schema column on a Table.
type Property struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	DataType string  `json:"dataType"`
}

// NewProperty returns a new Property.
func NewProperty(id int64, name string, typ string, dataType string) (*Property, error) {
  // Validate property type.
  if typ != ObjectType && typ != ActionType {
    return nil, fmt.Errorf("Invalid property type: %v", typ)
  }

  // Validate data type.
  switch dataType {
  case StringDataType, IntegerDataType, FloatDataType, BooleanDataType:
  default:
      return nil, fmt.Errorf("Invalid property data type: %v", dataType)
  }

  return &Property{
    Id: id,
    Name: name,
    Type: typ,
    DataType: dataType,
  }, nil
}

