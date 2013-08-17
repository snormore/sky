package query

import (
	"testing"
)

func TestStatementParserSelection(t *testing.T) {
    str := `SELECT count();`
    s := NewStatementParser().ParseString(&Query{}, str)
    if s.String() != str {
        t.Fatal("Unexpected:", "'"+s.String()+"'")
    }
}

func TestStatementParserCondition(t *testing.T) {
    str := `WHEN action == "signup" THEN`+"\n"+`  SELECT count();`+"\n"+`END`
    s := NewStatementParser().ParseString(&Query{}, str)
    if s.String() != str {
        t.Fatal("Unexpected:", "'"+s.String()+"'")
    }
}
