package skyd

import (
	"io/ioutil"
	"os"
	"time"
	"testing"
)

// Ensure that we can open and close a tablet.
func TestOpen(t *testing.T) {
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	
	tablet := NewTablet(path)
	defer tablet.Close()
	err = tablet.Open()
	if err != nil {
		t.Fatalf("Unable to open tablet: %v", err)
  }
}

// Ensure that we can add events and read them back.
func TestAddEvent(t *testing.T) {
	// Setup blank database.
	path, err := ioutil.TempDir("", "")
	defer os.RemoveAll(path)
	tablet := NewTablet(path)
	defer tablet.Close()
	_ = tablet.Open()

  // Setup source events.
  input := make([]*Event, 1)
  input[0] = &Event{}
  input[0].Timestamp, _ = time.Parse(time.RFC3339, "2012-01-01T00:00:00Z")
  input[0].Action = map[int64]interface{}{1:"foo"}
  input[0].Data = map[int64]interface{}{}

  err = tablet.AddEvent("bob", input[0])
	if err != nil {
		t.Fatalf("Unable to add event: %v", err)
  }
  
  // Read events out.
  output, err := tablet.GetEvents("bob")
  if err != nil {
    t.Fatalf("Unable to retrieve events: %v", err)
  }
  if len(output) != len(input) {
    t.Fatalf("Expected %v events, received %v", len(input), len(output))
  }
  for i := range output {
    if !input[i].Equal(output[i]) {
      t.Fatalf("Events not equal:\n  IN:  %v\n  OUT: %v", input[i], output[i])
    }
  }
}

