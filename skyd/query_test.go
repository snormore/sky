package skyd

import (
	"bytes"
	"testing"
)

// Ensure that we can encode and decode queries.
func TestQueryEncode(t *testing.T) {
	table := createTempTable(t)
	table.Open()
	defer table.Close()

	q := NewQuery(table)

	s1 := NewQueryCondition(q)
	s1.Within = 2
	s1.WithinUnits = "steps"
	q.Steps = append(q.Steps, s1)
	
	s2 := NewQuerySelection(q)
	s2.Expression = "count()"
	s2.Alias = "count"
	s2.Dimensions = []string{"foo", "bar"}
	q.Steps = append(q.Steps, s2)

	s3 := NewQuerySelection(q)
	s3.Expression = "sum(x)"
	s3.Alias = "myValue"
	s3.Dimensions = []string{}
	s1.Steps = append(s1.Steps, s3)
	
	buffer := new(bytes.Buffer)
	q.Encode(buffer)
	expected := `{"steps":[{"steps":[{"alias":"myValue","dimensions":[],"expression":"sum(x)","steps":[],"type":"selection"}],"type":"condition","within":2,"withinUnits":"steps"},{"alias":"count","dimensions":["foo","bar"],"expression":"count()","steps":[],"type":"selection"}]}`+"\n"
	if buffer.String() != expected {
		t.Fatalf("Query encoding error:\nexp: %s\ngot: %s", expected, buffer.String())
	}
}
