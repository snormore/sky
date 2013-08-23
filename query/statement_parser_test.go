package query

import (
	"testing"
)

func TestStatementParserSelection(t *testing.T) {
	str := `SELECT count() AS count, sum(price) AS price_sum`
	s, err := NewStatementParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if s.String() != str {
		t.Fatal("Unexpected:", "'"+s.String()+"'")
	}
}

func TestStatementParserCondition(t *testing.T) {
	str := `WHEN action == "signup" THEN` + "\n" + `  SELECT count() AS count` + "\n" + `END`
	s, err := NewStatementParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if s.String() != str {
		t.Fatal("Unexpected:", "'"+s.String()+"'")
	}
}
