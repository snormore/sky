package query

import (
	"encoding/json"
	"testing"
)

func TestParserSelection(t *testing.T) {
	str := `SELECT xyz;`
	query := NewParser().ParseString(str)
	expected := ``
	if b, _ := json.Marshal(query); string(b) != expected {
		t.Fatal("Unexpected:", string(b))
	}
}
