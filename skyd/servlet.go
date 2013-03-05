package skyd

import (
  "bytes"
  "errors"
  "fmt"
  "github.com/jmhodges/levigo"
  "io"
  "sort"
  "os"
)

// A Servlet is a small wrapper around a single shard of a LevelDB data file.
type Servlet struct {
  path string
  db   *levigo.DB
}

// NewServlet returns a new Servlet with a data shard stored at a given path.
func NewServlet(path string) *Servlet {
  return &Servlet{
    path: path,
  }
}

// Opens the underlying LevelDB database.
func (s *Servlet) Open() error {
  err := os.MkdirAll(s.path, 0700)
  if err != nil {
    return err
  }

  opts := levigo.NewOptions()
  opts.SetCreateIfMissing(true)
  db, err := levigo.Open(s.path, opts)
  if err != nil {
    panic(fmt.Sprintf("skyd.Servlet: Unable to open LevelDB database: %v", err))
  }
  s.db = db
  return nil
}

// Closes the underlying LevelDB database.
func (s *Servlet) Close() {
  if s.db != nil {
    s.db.Close()
  }
}

// Adds an event for a given object in a table to a servlet.
func (s *Servlet) PutEvent(table *Table, objectId string, event *Event) error {
  // Make sure the servlet is open.
  if s.db == nil {
    return fmt.Errorf("Servlet is not open: %v", s.path)
  }

  // Do not allow empty events to be added.
  if event == nil {
    return errors.New("skyd.PutEvent: Cannot add nil event")
  }

  // Retrieve the events for the object and append.
  events, err := s.GetEvents(table, objectId)
  if err != nil {
    return err
  }
  events = append(events, event)

  // Write events back to the database.
  err = s.SetEvents(table, objectId, events)
  if err != nil {
    return err
  }

  return nil
}

// Retrieves a list of events for a given object in a table.
func (s *Servlet) GetEvents(table *Table, objectId string) ([]*Event, error) {
  // Make sure the servlet is open.
  if s.db == nil {
    return nil, fmt.Errorf("Servlet is not open: %v", s.path)
  }

  // Encode object identifier.
  encodedObjectId, err := table.EncodeObjectId(objectId)
  if err != nil {
    return nil, err
  }

  // Retrieve byte array.
  ro := levigo.NewReadOptions()
  data, err := s.db.Get(ro, encodedObjectId)
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

// Writes a list of events for an object in table.
func (s *Servlet) SetEvents(table *Table, objectId string, events []*Event) error {
  // Make sure the servlet is open.
  if s.db == nil {
    return fmt.Errorf("Servlet is not open: %v", s.path)
  }

  // Sort the events.
  sort.Sort(EventList(events))

  // Encode object identifier.
  encodedObjectId, err := table.EncodeObjectId(objectId)
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
  err = s.db.Put(wo, encodedObjectId, buffer.Bytes())
  wo.Close()

  return nil
}
