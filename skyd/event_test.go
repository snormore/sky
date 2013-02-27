package skyd

import (
	"bytes"
	"time"
	"testing"
)

// Encode and then decode an Event to MsgPack.
func TestEncodeDecode(t *testing.T) {
  timeString := "1970-01-01T00:01:00Z"
  timestamp, err := time.Parse(time.RFC3339, timeString)
	e1 := &Event{Timestamp:timestamp}
	e1.Action = map[int64]interface{}{1:int64(20),-2:"baz"}
	e1.Data = map[int64]interface{}{10:int64(123)}

  // Encode
	buffer := new(bytes.Buffer)
  err = e1.EncodeRaw(buffer)
	if err != nil {
		t.Errorf("Unable to encode: %v", err)
  }

  // Decode
  e2 := &Event{}
  err = e2.DecodeRaw(buffer)
  if err != nil {
    t.Errorf("Unable to decode: %v", err)
  }
  if !e1.Equal(e2) {
    t.Errorf("Events do not match: %v <=> %v", e1, e2)
  }
}

