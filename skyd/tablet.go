package skyd

import (
	"github.com/jmhodges/levigo"
	"bytes"
	"errors"
	"fmt"
	"io"
	"sort"
)

// A Tablet is a small wrapper for the underlying data storage
// which is contained in LevelDB.
type Tablet struct {
	db   *levigo.DB
	Path string
}

// NewTablet returns a new Tablet that is stored at a given path.
func NewTablet(path string) Tablet {
  return Tablet{
    Path: path,
  }
}

// Opens the underlying LevelDB database.
func (t *Tablet) Open() error {
	opts := levigo.NewOptions()
	opts.SetCreateIfMissing(true)
	db, err := levigo.Open(t.Path, opts)
	if err != nil {
		return err
	}
	t.db = db
	return nil
}

// Closes the underlying LevelDB database.
func (t *Tablet) Close() {
  if t.db != nil {
	  t.db.Close()
  }
}


// Adds an event for a given object to a tablet.
func (t *Tablet) AddEvent(objectId interface{}, event *Event) error {
  // Make sure the tablet is open.
  if t.db == nil {
	  return fmt.Errorf("Tablet is not open: %v", t.Path)
  }

  // Do not allow empty events to be added.
  if event == nil {
    return errors.New("skyd.AddEvent: Cannot add nil event")
  }
  
  // Retrieve the events for the object and append.
  events, err := t.GetEvents(objectId)
  if err != nil {
    return err
  }
  events = append(events, event)
  
  // Write events back to the database.
  err = t.SetEvents(objectId, events)
  if err != nil {
    return err
  }
  
  return nil
}

// Retrieves a list of events for a given object.
func (t *Tablet) GetEvents(objectId interface{}) ([]*Event, error) {
  // Make sure the tablet is open.
  if t.db == nil {
	  return nil, fmt.Errorf("Tablet is not open: %v", t.Path)
  }
  
  // Encode object identifier.
  encodedObjectId, err := EncodeObjectId(objectId)
  if err != nil {
    return nil, err
  }
  
  // Retrieve byte array.
  ro := levigo.NewReadOptions()
  data, err := t.db.Get(ro, encodedObjectId)
  ro.Close()
  if err != nil {
    return nil, err
  }
  
  // Decode the events into a slice.
  events := make([]*Event, 0, 10)
  if data != nil {
    reader := bytes.NewReader(data)
    for {
      // Otherwise decode the event and append it to our list.
      event := &Event{}
      err = event.DecodeRaw(reader)
      if err == io.EOF {
        break
      }
      if err != nil {
        return nil, err
      }
      events = append(events, event)
    }
  }
  
  return events, nil
}

// Writes a list of events to the database.
func (t *Tablet) SetEvents(objectId interface{}, events []*Event) error {
  // Make sure the tablet is open.
  if t.db == nil {
	  return fmt.Errorf("Tablet is not open: %v", t.Path)
  }

  // Sort the events.
  sort.Sort(EventSlice(events))

  // Encode object identifier.
  encodedObjectId, err := EncodeObjectId(objectId)
  if err != nil {
    return err
  }

  // Encode the events.
  buffer := new(bytes.Buffer)
  for _, event := range events {
    err = event.EncodeRaw(buffer)
    if err != nil {
      return err
    }
  }

  // Write bytes to the database.
  wo := levigo.NewWriteOptions()
  err = t.db.Put(wo, encodedObjectId, buffer.Bytes())
  wo.Close()
  
  return nil
}

