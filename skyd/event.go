package skyd

import (
	"fmt"
	"github.com/ugorji/go-msgpack"
	"io"
	"time"
)

// An Event is a state change that occurs at a particular point in time.
type Event struct {
	Timestamp time.Time
	Data      map[int64]interface{}
}

// NewEvent returns a new Event.
func NewEvent(timestamp string, data map[int64]interface{}) *Event {
	t, _ := time.Parse(time.RFC3339, timestamp)
	return &Event{
		Timestamp: t,
		Data:      data,
	}
}

// Encodes an event to MsgPack format.
func (e *Event) EncodeRaw(writer io.Writer) error {
	raw := []interface{}{ShiftTime(e.Timestamp), e.Data}
	encoder := msgpack.NewEncoder(writer)
	err := encoder.Encode(raw)
	return err
}

// Decodes an event from MsgPack format.
func (e *Event) DecodeRaw(reader io.Reader) error {
	raw := make([]interface{}, 2)
	decoder := msgpack.NewDecoder(reader, nil)
	err := decoder.Decode(&raw)
	if err != nil {
		return err
	}

	// Convert the timestamp to int64.
	timestamp, err := castInt64(raw[0])
	if err != nil {
		return fmt.Errorf("Unable to parse timestamp: '%v'", raw[0])
	}
	e.Timestamp = UnshiftTime(timestamp).UTC()

	// Convert data to appropriate map.
	if raw[1] != nil {
		e.Data, err = e.decodeRawMap(raw[1].(map[interface{}]interface{}))
		if err != nil {
			return err
		}
	}

	return nil
}

// Decodes the map.
func (e *Event) decodeRawMap(raw map[interface{}]interface{}) (map[int64]interface{}, error) {
	m := make(map[int64]interface{})
	for k, v := range raw {
		kInt64, err := castInt64(k)
		if err != nil {
			return nil, err
		}

		vInt64, err := castInt64(v)
		if err == nil {
			m[kInt64] = vInt64
		} else {
			m[kInt64] = v
		}
	}
	return m, nil
}

// Compares two events for equality.
func (e *Event) Equal(x *Event) bool {
	if !e.Timestamp.Equal(x.Timestamp) {
		return false
	}
	for k, v := range e.Data {
		v2 := x.Data[k]
		if v != v2 {
			return false
		}
	}
	for k, v := range x.Data {
		v2 := e.Data[k]
		if v != v2 {
			return false
		}
	}
	return true
}

// Merges the data of another event into this event.
func (e *Event) Merge(a *Event) {
	if e.Data == nil && a.Data != nil {
		e.Data = make(map[int64]interface{})
	}
	for k, v := range a.Data {
		e.Data[k] = v
	}
}
