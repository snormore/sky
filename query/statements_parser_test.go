package query

import (
	"testing"
)

func TestStatementsParser(t *testing.T) {
	str := `WHEN action == "signup" THEN` + "\n" + `  SELECT count() AS count;` + "\n" + `END` + "\n" + `SELECT count() AS count, sum(price) AS price_sum;`
	s, err := NewStatementsParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if s.String() != str {
		t.Fatal("Unexpected:", "'"+s.String()+"'")
	}
}
