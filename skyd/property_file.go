package skyd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
)

// A PropertyFile manages the serialization of Property objects for a table.
type PropertyFile struct {
	path             string
	properties       map[int64]*Property
	propertiesByName map[string]*Property
}

// NewProperty returns a new PropertyFile.
func NewPropertyFile(path string) *PropertyFile {
	p := &PropertyFile{
		path: path,
	}
	p.Reset()
	return p
}

// Adds a new property to the property file and generate an identifier for it.
func (p *PropertyFile) CreateProperty(name string, transient bool, dataType string) (*Property, error) {
	// Don't allow duplicate names.
	if p.propertiesByName[name] != nil {
		return nil, fmt.Errorf("Property already exists: %v", name)
	}

	property, err := NewProperty(0, name, transient, dataType)
	if err != nil {
		return nil, err
	}

	// Find the next object/action identifier.
	if property.Transient {
		_, property.Id = p.NextIdentifiers()
	} else {
		property.Id, _ = p.NextIdentifiers()
	}

	// Add to the list.
	p.properties[property.Id] = property
	p.propertiesByName[property.Name] = property

	return property, nil
}

// Finds the next available action and object property identifiers.
func (p *PropertyFile) NextIdentifiers() (int64, int64) {
	var nextPermanentId, nextTransientId int64 = 1, -1
	for _, property := range p.properties {
		if property.Transient && property.Id <= nextTransientId {
			nextTransientId = property.Id - 1
		} else if !property.Transient && property.Id >= nextPermanentId {
			nextPermanentId = property.Id + 1
		}
	}
	return nextPermanentId, nextTransientId
}

// Encodes a property file.
func (p *PropertyFile) Encode(writer io.Writer) error {
	// Convert the lookup into a sorted slice.
	list := p.GetAllProperties()

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
	p.Reset()
	for _, property := range list {
		p.properties[property.Id] = property
		if property.Name != "" {
			p.propertiesByName[property.Name] = property
		}
	}

	return nil
}

// Saves the property file to disk.
func (p *PropertyFile) Save() error {
	// Open the file for writing.
	file, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Then decode it.
	w := bufio.NewWriter(file)
	err = p.Encode(w)
	if err != nil {
		return err
	}
	if err = w.Flush(); err != nil {
		return err
	}

	return nil
}

// Loads the property file from disk.
func (p *PropertyFile) Load() error {
	// Ignore if there is no file.
	if _, err := os.Stat(p.path); os.IsNotExist(err) {
		return nil
	}

	// Otherwise open it.
	file, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Then decode it.
	err = p.Decode(bufio.NewReader(file))
	if err != nil {
		return err
	}

	return nil
}

// Clears out the property file.
func (p *PropertyFile) Reset() {
	p.properties = make(map[int64]*Property)
	p.propertiesByName = make(map[string]*Property)
}

// Retrieves a list of undeleted properties sorted by id.
func (p *PropertyFile) GetProperties() []*Property {
	list := make([]*Property, 0)
	for _, property := range p.propertiesByName {
		list = append(list, property)
	}
	sort.Sort(PropertyList(list))
	return list
}

// Retrieves a list of all properties sorted by id.
func (p *PropertyFile) GetAllProperties() []*Property {
	list := make([]*Property, 0)
	for _, property := range p.properties {
		list = append(list, property)
	}
	sort.Sort(PropertyList(list))
	return list
}

// Retrieves a single property by id.
func (p *PropertyFile) GetProperty(id int64) *Property {
	return p.properties[id]
}

// Retrieves a single property by name.
func (p *PropertyFile) GetPropertyByName(name string) *Property {
	return p.propertiesByName[name]
}

// Deletes a property.
func (p *PropertyFile) DeleteProperty(property *Property) {
	if property != nil && property.Name != "" {
		delete(p.propertiesByName, property.Name)
		property.Name = ""
	}
}

// Converts a map with string keys to use property identifier keys.
func (p *PropertyFile) NormalizeMap(m map[string]interface{}) (map[int64]interface{}, error) {
	clone := make(map[int64]interface{})
	for k, v := range m {
		// Look up the property by name and convert it to the ID.
		property := p.GetPropertyByName(string(k))
		if property != nil {
			clone[property.Id] = v
		} else {
			return nil, fmt.Errorf("Property not found: %v", k)
		}
	}
	return clone, nil
}

// Converts a map with property identifier keys to use string keys.
func (p *PropertyFile) DenormalizeMap(m map[int64]interface{}) (map[string]interface{}, error) {
	clone := make(map[string]interface{})
	for k, v := range m {
		// Look up the property by ID and convert it to the name.
		property := p.GetProperty(k)
		if property != nil {
			clone[property.Name] = v
		} else {
			return nil, fmt.Errorf("skyd.PropertyFile: Property not found: %v", k)
		}
	}
	return clone, nil
}
