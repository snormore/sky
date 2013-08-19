package query

import (
	"testing"
)

func TestExpressionParserSimple(t *testing.T) {
	str := `timestamp >= 2 && timestamp < 6`
	e, err := NewExpressionParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if e.String() != `(timestamp >= 2) && (timestamp < 6)` {
		t.Fatal("Unexpected:", "'"+e.String()+"'")
	}
}

func TestExpressionParserComplex(t *testing.T) {
	str := `x * 1 + (2 / 3) > 100`
	e, err := NewExpressionParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if e.String() != `((x * 1) + (2 / 3)) > 100` {
		t.Fatal("Unexpected:", "'"+e.String()+"'")
	}
}
