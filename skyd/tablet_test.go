package skyd

import (
  "io/ioutil"
  "os"
  "testing"
)

// Ensure that we can open and close a tablet.
func TestOpen(t *testing.T) {
  path, err := ioutil.TempDir("", "")
  defer os.RemoveAll(path)
  
  table := NewTable("/tmp/foo")
  tablet := NewTablet(table, path)
  defer tablet.Close()
  err = tablet.Open()
  if err != nil {
    t.Fatalf("Unable to open tablet: %v", err)
  }
}

// Ensure that we can add events and read them back.
func TestTabletAddEvent(t *testing.T) {
  // Setup blank database.
  path, err := ioutil.TempDir("", "")
  defer os.RemoveAll(path)
  table := NewTable("/tmp/foo")
  tablet := NewTablet(table, path)
  defer tablet.Close()
  _ = tablet.Open()

  // Setup source events.
  input := make([]*Event, 2)
  input[0] = NewEvent("2012-01-02T00:00:00Z", map[int64]interface{}{1:"foo"}, map[int64]interface{}{})
  input[1] = NewEvent("2012-01-01T00:00:00Z", map[int64]interface{}{2:"bar"}, map[int64]interface{}{})

  // Add events.
  for _, e := range input {
    err = tablet.AddEvent("bob", e)
    if err != nil {
      t.Fatalf("Unable to add event: %v", err)
    }
  }
  
  // Setup expected events.
  expected := make([]*Event, len(input))
  expected[0] = input[1]
  expected[1] = input[0]
  
  // Read events out.
  output, err := tablet.GetEvents("bob")
  if err != nil {
    t.Fatalf("Unable to retrieve events: %v", err)
  }
  if len(output) != len(expected) {
    t.Fatalf("Expected %v events, received %v", len(expected), len(output))
  }
  for i := range output {
    if !expected[i].Equal(output[i]) {
      t.Fatalf("Events not equal:\n  IN:  %v\n  OUT: %v", expected[i], output[i])
    }
  }
}

