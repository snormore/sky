package skyd

import (
  "errors"
  "fmt"
  "github.com/ugorji/go-msgpack"
  "os"
  "path/filepath"
  "time"
)

// A Table is a collection of objects.
type Table struct {
  name         string
  path         string
  propertyFile *PropertyFile
}

// NewTable returns a new Table that is stored at a given path.
func NewTable(name string, path string) *Table {
  path, err := filepath.Abs(path)
  if err != nil {
    return nil
  }

  return &Table{
    name:         name,
    path:         path,
  }
}

// Retrieves the path on the table.
func (t *Table) Path() string {
  return t.path
}

// Retrieves the name of the table.
func (t *Table) Name() string {
  return t.name
}

// Creates a table directory structure.
func (t *Table) Create() error {
  if t.Exists() {
    return fmt.Errorf("Table already exist: %v", t.name)
  }

  // Create root directory.
  err := os.MkdirAll(t.path, 0700)
  if err != nil {
    return err
  }

  return nil
}

// Deletes a table.
func (t *Table) Delete() error {
  if !t.Exists() {
    return fmt.Errorf("Table does not exist: %v", t.name)
  }

  // Close everything if it's open.
  if t.IsOpen() {
    t.Close()
  }

  // Delete the whole damn directory.
  os.RemoveAll(t.path)

  return nil
}

// Opens the table.
func (t *Table) Open() error {
  if !t.Exists() {
    return errors.New("Table does not exist")
  }

  // Load property file.
  t.propertyFile = NewPropertyFile(fmt.Sprintf("%v/%v", t.path, "properties"))
  err := t.propertyFile.Load()
  if err != nil {
    t.Close()
    return err
  }

  return nil
}

// Closes the table.
func (t *Table) Close() {
  t.propertyFile = nil
}

// Checks if the table is currently open.
func (t *Table) IsOpen() bool {
  return t.propertyFile != nil
}

// Checks if the table exists on disk.
func (t *Table) Exists() bool {
  if _, err := os.Stat(t.path); os.IsNotExist(err) {
    return false
  }
  return true
}

// Adds a property to the table.
func (t *Table) CreateProperty(name string, typ string, dataType string) (*Property, error) {
  if !t.IsOpen() {
    return nil, errors.New("Table is not open")
  }

  // Create property on property file.
  property, err := t.propertyFile.CreateProperty(name, typ, dataType)
  if err != nil {
    return nil, err
  }

  // Save the property file to disk.
  err = t.propertyFile.Save()
  if err != nil {
    return nil, err
  }

  return property, err
}

// Retrieves a list of all properties on the table.
func (t *Table) GetProperties() ([]*Property, error) {
  if !t.IsOpen() {
    return nil, errors.New("Table is not open")
  }
  return t.propertyFile.GetProperties(), nil
}

// Retrieves a single property from the table by id.
func (t *Table) GetProperty(id int64) (*Property, error) {
  if !t.IsOpen() {
    return nil, errors.New("Table is not open")
  }
  return t.propertyFile.GetProperty(id), nil
}

// Retrieves a single property from the table by name.
func (t *Table) GetPropertyByName(name string) (*Property, error) {
  if !t.IsOpen() {
    return nil, errors.New("Table is not open")
  }
  return t.propertyFile.GetPropertyByName(name), nil
}

// Deletes a single property on the table.
func (t *Table) DeleteProperty(property *Property) error {
  if !t.IsOpen() {
    return errors.New("Table is not open")
  }
  t.propertyFile.DeleteProperty(property)
  return nil
}

// Saves the property file on the table.
func (t *Table) SavePropertyFile() error {
  if !t.IsOpen() {
    return errors.New("Table is not open")
  }
  return t.propertyFile.Save()
}

// Converts a map with string keys to use property identifier keys.
func (t *Table) NormalizeMap(m map[string]interface{}) (map[int64]interface{}, error) {
  return t.propertyFile.NormalizeMap(m)
}

// Converts a map with property identifier keys to use string keys.
func (t *Table) DenormalizeMap(m map[int64]interface{}) (map[string]interface{}, error) {
  return t.propertyFile.DenormalizeMap(m)
}

// Encodes an object identifier for this table.
func (t *Table) EncodeObjectId(objectId string) ([]byte, error) {
  return msgpack.Marshal([]string{t.name,objectId})
}


// Deserializes a map into a normalized event.
func (t *Table) DeserializeEvent(m map[string]interface{}) (*Event, error) {
  event := &Event{}
  
  // Parse timestamp.
  if timestamp, ok := m["timestamp"].(string); ok {
    ts, err := time.Parse(time.RFC3339, timestamp)
    if err != nil {
      return nil, fmt.Errorf("Unable to parse timestamp: %v", timestamp)
    }
    event.Timestamp = ts
  } else {
    return nil, errors.New("Timestamp required.")
  }
  
  // Convert maps to use property identifiers.
  if data, ok := m["data"].(map[string]interface{}); ok {
    normalizedData, err := t.NormalizeMap(data)
    if err != nil {
      return nil, err
    }
    event.Data = normalizedData
  }
  
  return event, nil
}

// Serializes a normalized event into a map.
func (t *Table) SerializeEvent(event *Event) (map[string]interface{}, error) {
  m := make(map[string]interface{})
  
  // Format timestamp.
  m["timestamp"] = event.Timestamp.UTC().Format(time.RFC3339)
  
  // Convert data map to use property names.
  if event.Data != nil {
    denormalizedData, err := t.DenormalizeMap(event.Data)
    if err != nil {
      return nil, err
    }
    m["data"] = denormalizedData
  }

  return m, nil
}
