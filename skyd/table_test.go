package skyd

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

// Ensure that we can create a new table.
func TestCreate(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	path = fmt.Sprintf("%v/test", path)

	table := NewTable(path)
	err = table.Create()
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}
	if !table.Exists() {
		t.Fatalf("Table doesn't exist: %v", path)
	}
	for i := 0; i < runtime.NumCPU(); i++ {
		if _, err := os.Stat(fmt.Sprintf("%v/%v", path, i)); os.IsNotExist(err) {
			t.Fatalf("Tablet directory doesn't exist: %v", i)
		}
	}
}

// Ensure that we can add an event to the appropriate tablet.
func TestTableAddEvent(t *testing.T) {
	table := createTable(t)
	table.Open()
	defer table.Close()

	input := NewEvent("2012-01-02T00:00:00Z", map[int64]interface{}{1: "foo"}, map[int64]interface{}{})
	err := table.AddEvent("bob", input)
	if err != nil {
		t.Fatalf("Unable to add event to table: %v", err)
	}

	output, err := table.GetEvents("bob")
	if err != nil {
		t.Fatalf("Unable to retrieve events: %v", err)
	}
	if len(output) != 1 {
		t.Fatalf("Expected %v events, received %v", 1, len(output))
	}
	for i := range output {
		if !input.Equal(output[i]) {
			t.Fatalf("Events not equal:\n  IN:  %v\n  OUT: %v", input, output[i])
		}
	}
}

// Ensure that we can create properties on a table.
func TestTableCreateProperty(t *testing.T) {
	table := createTable(t)
	table.Open()
	defer table.Close()

	property, err := table.CreateProperty("name", "object", "string")
	if property == nil || err != nil {
		t.Fatalf("Unable to add property to table: %v", err)
	}

  content, _ := ioutil.ReadFile(fmt.Sprintf("%v/properties", table.Path()))
  if string(content) != "[{\"id\":1,\"name\":\"name\",\"type\":\"object\",\"dataType\":\"string\"}]\n" {
    t.Fatalf("Invalid properties file:\n%v", string(content))
  }
}

// Creates a table.
func createTable(t *testing.T) *Table {
	path, err := ioutil.TempDir("", "")
	os.RemoveAll(path)

	table := NewTable(path)
	err = table.Create()
	if err != nil {
		t.Fatalf("Unable to create table: %v", err)
	}

	return table
}
