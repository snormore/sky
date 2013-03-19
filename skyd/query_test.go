package skyd

import (
	"bytes"
	"testing"
)

// Ensure that we can encode queries.
func TestQueryEncodeDecode(t *testing.T) {
	table := createTempTable(t)
	table.Open()
	defer table.Close()

	json := `{"sessionIdleTime":0,"steps":[{"expression":"baz == 'hello'","steps":[{"alias":"myValue","dimensions":[],"expression":"sum(x)","steps":[],"type":"selection"}],"type":"condition","within":2,"withinUnits":"steps"},{"alias":"count","dimensions":["foo","bar"],"expression":"count()","steps":[],"type":"selection"}]}` + "\n"

	// Decode
	q := NewQuery(table)
	buffer := bytes.NewBufferString(json)
	err := q.Decode(buffer)
	if err != nil {
		t.Fatalf("Query decoding error: %v", err)
	}

	// Encode
	buffer = new(bytes.Buffer)
	q.Encode(buffer)
	if buffer.String() != json {
		t.Fatalf("Query encoding error:\nexp: %s\ngot: %s", json, buffer.String())
	}
}

