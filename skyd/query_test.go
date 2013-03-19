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

// Ensure that we can codegen queries.
func TestQueryCodegen(t *testing.T) {
	table := createTempTable(t)
	table.Open()
	defer table.Close()

	q := NewQuery(table)
	err := q.Decode(bytes.NewBufferString(`{
		"steps":[
			{"type":"condition","expression":"bar == 'baz'","within":0,"withinUnits":"steps","steps":[
				{"type":"selection","alias":"foo","dimensions":["xxx","yyy"],"expression":"count()","steps":[]}
			]}
		]
	}`))
	if err != nil {
		t.Fatalf("Query decoding error: %v", err)
	}

	// Codegen
	code, err := q.Codegen()
	if err != nil {
		t.Fatalf("Query codegen error: %v", err)
	}
	exp := `_`
	if code != exp {
		t.Fatalf("Query codegen:\nexp:\n%s\ngot:\n%s", code, exp)
	}
}
