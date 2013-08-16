package query

import (
	"encoding/json"
	"testing"
)

func TestParserSelectCount(t *testing.T) {
	str := `SELECT count();`
	query := NewParser().ParseString(str)
	expected := ``
	if b, _ := json.Marshal(query); string(b) != expected {
		t.Fatal("Unexpected:", string(b))
	}
}
