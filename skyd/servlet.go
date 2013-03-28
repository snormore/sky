package skyd

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jmhodges/levigo"
	"io"
	"os"
	"sort"
	"time"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A Servlet is a small wrapper around a single shard of a LevelDB data file.
type Servlet struct {
	path    string
	db      *levigo.DB
	factors *Factors
	channel chan *Message
}

//------------------------------------------------------------------------------
//
// Constructors
//
//------------------------------------------------------------------------------

// NewServlet returns a new Servlet with a data shard stored at a given path.
func NewServlet(path string, factors *Factors) *Servlet {
	return &Servlet{
		path:    path,
		factors: factors,
		channel: make(chan *Message),
	}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Lifecycle
//--------------------------------------

// Opens the underlying LevelDB database and starts the message loop.
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

	// Start message loop.
	go s.process()

	return nil
}

// Closes the underlying LevelDB database.
func (s *Servlet) Close() {
	// Wait for the message loop to shutdown.
	m := NewShutdownMessage()
	s.channel <- m
	m.wait()

	// Cleanup channel.
	close(s.channel)
	s.channel = nil

	if s.db != nil {
		s.db.Close()
	}
}

//--------------------------------------
// Synchronization
//--------------------------------------

// Executes a function through a single-threaded servlet context and waits until completion.
func (s *Servlet) sync(f func() (interface{}, error)) (interface{}, error) {
	m := s.async(f)
	return m.wait()
}

// Executes a function through a single-threaded servlet context and returns immediately.
func (s *Servlet) async(f func() (interface{}, error)) *Message {
	m := NewExecuteMessage(f)
	s.channel <- m
	return m
}

//--------------------------------------
// Message Processing
//--------------------------------------

// Serially processes messages routed through the servlet channel.
func (s *Servlet) process() {
	for message := range s.channel {
		message.execute()
		if message.messageType == ShutdownMessageType {
			return
		}
	}
}

//--------------------------------------
// Event Management
//--------------------------------------

// Adds an event for a given object in a table to a servlet.
func (s *Servlet) PutEvent(table *Table, objectId string, event *Event, replace bool) error {
	// Make sure the servlet is open.
	if s.db == nil {
		return fmt.Errorf("Servlet is not open: %v", s.path)
	}

	// Do not allow empty events to be added.
	if event == nil {
		return errors.New("skyd.PutEvent: Cannot add nil event")
	}

	// Retrieve the events for the object and append.
	tmp, err := s.GetEvents(table, objectId)
	if err != nil {
		return err
	}

	// Remove any event matching the timestamp.
	found := false
	state := &Event{Timestamp:event.Timestamp, Data:map[int64]interface{}{}}
	events := make([]*Event, 0)
	for _, v := range tmp {
		// Replace or merge with existing event.
		if v.Timestamp.Equal(event.Timestamp) {
			// Dedupe all permanent state.
			event.Dedupe(state)
			
			// Replace or merge.
			if replace {
				v = event
			} else {
				v.Merge(event)
			}
			found = true
		}
		events = append(events, v)

		// Keep track of permanent state.
		state.MergePermanent(v)
	}
	// Add the event if it wasn't found.
	if !found {
		event.Dedupe(state)
		events = append(events, event)
	}

	// Write events back to the database.
	err = s.SetEvents(table, objectId, events)
	if err != nil {
		return err
	}

	return nil
}

// Retrieves an event for a given object at a single point in time.
func (s *Servlet) GetEvent(table *Table, objectId string, timestamp time.Time) (*Event, error) {
	// Retrieve all events.
	events, err := s.GetEvents(table, objectId)
	if err != nil {
		return nil, err
	}

	// Find an event at a given point in time.
	for _, v := range events {
		if v.Timestamp.Equal(timestamp) {
			return v, nil
		}
	}

	return nil, nil
}

// Removes an event for a given object in a table to a servlet.
func (s *Servlet) DeleteEvent(table *Table, objectId string, timestamp time.Time) error {
	// Make sure the servlet is open.
	if s.db == nil {
		return fmt.Errorf("Servlet is not open: %v", s.path)
	}

	// Retrieve the events for the object and append.
	tmp, err := s.GetEvents(table, objectId)
	if err != nil {
		return err
	}
	// Remove any event matching the timestamp.
	events := make([]*Event, 0)
	for _, v := range tmp {
		if !v.Timestamp.Equal(timestamp) {
			events = append(events, v)
		}
	}

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

// Deletes all events for a given object in a table.
func (s *Servlet) DeleteEvents(table *Table, objectId string) error {
	// Make sure the servlet is open.
	if s.db == nil {
		return fmt.Errorf("Servlet is not open: %v", s.path)
	}

	// Encode object identifier.
	encodedObjectId, err := table.EncodeObjectId(objectId)
	if err != nil {
		return err
	}

	// Delete object from the database.
	wo := levigo.NewWriteOptions()
	err = s.db.Delete(wo, encodedObjectId)
	wo.Close()

	return nil
}
