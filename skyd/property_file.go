package skyd

import (
  "encoding/json"
  "fmt"
  "io"
  "sort"
)

// A PropertyFile manages the serialization of Property objects for a table.
type PropertyFile struct {
	properties   map[int64]*Property
	propertiesByName map[string]*Property
}

// NewProperty returns a new PropertyFile.
func NewPropertyFile() *PropertyFile {
	return &PropertyFile{
	  properties: make(map[int64]*Property),
	  propertiesByName: make(map[string]*Property),
	}
}

// Adds a new property to the property file and generate an identifier for it.
func (p *PropertyFile) CreateProperty(name string, typ string, dataType string) (*Property, error) {
  // Don't allow duplicate names.
  if p.propertiesByName[name] != nil {
    return nil, fmt.Errorf("Property already exists: %v", name)
  }
  
  property, err := NewProperty(0, name, typ, dataType)
  if err != nil {
    return nil, err
  }
  
  // Find the next object/action identifier.
  if property.Type == ObjectType {
    property.Id, _ = p.NextIdentifiers()
  } else {
    _, property.Id = p.NextIdentifiers()
  }

  // Add to the list.
  p.properties[property.Id] = property
  p.propertiesByName[property.Name] = property

  return property, nil
}

// Finds the next available action and object property identifiers.
func (p *PropertyFile) NextIdentifiers() (int64, int64) {
  var nextObjectId, nextActionId int64 = 1, -1
  for _, property := range p.properties {
    if property.Type == ObjectType && property.Id >= nextObjectId {
      nextObjectId = property.Id + 1
    } else if property.Type == ActionType && property.Id <= nextActionId {
      nextActionId = property.Id - 1
    }
  }
  return nextObjectId, nextActionId
}

// Encodes a property file.
func (p *PropertyFile) Encode(writer io.Writer) error {
	// Convert the lookup into a sorted slice.
	list := make([]*Property, 0)
	for _, property := range p.properties {
	  list = append(list, property)
	}
  sort.Sort(PropertyList(list))

  // Encode the slice.
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(list)
	return err
}

// Decodes a property file.
func (p *PropertyFile) Decode(reader io.Reader) error {
  list := make([]*Property, 0)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&list)
	if err != nil {
	  return err
	}
	
	// Create lookups for the properties.
	p.properties = make(map[int64]*Property)
	p.propertiesByName = make(map[string]*Property)
	for _, property := range list {
	  p.properties[property.Id] = property
	  p.propertiesByName[property.Name] = property
	}
	
	return nil
}

