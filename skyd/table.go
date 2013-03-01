package skyd

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

// A Table is a collection of tablets.
type Table struct {
	name         string
	path         string
	tablets      []*Tablet
	propertyFile *PropertyFile
}

// NewTable returns a new Table that is stored at a given path.
func NewTable(path string) *Table {
	path, err := filepath.Abs(path)
	if err != nil {
		return nil
	}

	return &Table{
		path:         path,
		name:         filepath.Base(path),
		propertyFile: NewPropertyFile(fmt.Sprintf("%v/%v", path, "properties")),
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

	// Create a subdirectory for each tablet.
	for i := 0; i < runtime.NumCPU(); i++ {
		err = os.Mkdir(fmt.Sprintf("%v/%v", t.path, i), 0700)
		if err != nil {
			os.RemoveAll(t.path)
			return nil
		}
	}

	return nil
}

// Opens the table.
func (t *Table) Open() error {
	if !t.Exists() {
		return errors.New("Table does not exist")
	}

  // Load property file.
  err := t.propertyFile.Load()
  if err != nil {
    t.Close()
    return err
  }

	// Create tablets from child directories with numeric names.
	infos, err := ioutil.ReadDir(t.path)
	if err != nil {
		return err
	}
	for _, info := range infos {
		match, _ := regexp.MatchString("^\\d$", info.Name())
		if info.IsDir() && match {
			tablet := NewTablet(fmt.Sprintf("%s/%s", t.path, info.Name()))
			t.tablets = append(t.tablets, tablet)
			err = tablet.Open()
			if err != nil {
	      t.Close()
        return err
			}
		}
	}

	return nil
}

// Closes the table and all the tablets.
func (t *Table) Close() {
  t.propertyFile.Reset()

	for _, tablet := range t.tablets {
		tablet.Close()
	}
	t.tablets = nil
}

// Checks if the table is currently open.
func (t *Table) IsOpen() bool {
	return t.tablets != nil
}

// Checks if the table exists on disk.
func (t *Table) Exists() bool {
	if _, err := os.Stat(t.path); os.IsNotExist(err) {
		return false
	}
	return true
}

// Adds an event for a given object to the table.
func (t *Table) AddEvent(objectId interface{}, event *Event) error {
	if !t.IsOpen() {
		return errors.New("Table is not open")
	}

	// Determine tablet number that event should go to.
	tabletIndex, err := t.GetObjectTabletIndex(objectId)
	if err != nil {
		return err
	}

	// Add event to the appropriate tablet.
	tablet := t.tablets[tabletIndex]
	err = tablet.AddEvent(objectId, event)
	if err != nil {
		return err
	}

	return nil
}

// Retrieves a list of events for a given object.
func (t *Table) GetEvents(objectId interface{}) ([]*Event, error) {
	if !t.IsOpen() {
		return nil, errors.New("Table is not open")
	}

	// Determine tablet number that event should go to.
	tabletIndex, err := t.GetObjectTabletIndex(objectId)
	if err != nil {
		return nil, err
	}

	// Add event to the appropriate tablet.
	tablet := t.tablets[tabletIndex]
	events, err := tablet.GetEvents(objectId)
	if err != nil {
		return nil, err
	}

	return events, nil
}

// Calculates a tablet index based on the object identifier even hash.
func (t *Table) GetObjectTabletIndex(objectId interface{}) (uint32, error) {
	if !t.IsOpen() {
		return 0, errors.New("Table is not open")
	}

	// Encode object identifier.
	encodedObjectId, err := EncodeObjectId(objectId)
	if err != nil {
		return 0, err
	}

	// Calculate the even bits of the FNV1a hash.
	h := fnv.New64a()
	h.Reset()
	h.Write(encodedObjectId)
	hashcode := h.Sum64()
	index := CondenseUint64Even(hashcode) % uint32(len(t.tablets))

	return index, nil
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
