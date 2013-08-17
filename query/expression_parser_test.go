package query

import (
	"testing"
)

func TestExpressionParserComplex(t *testing.T) {
	str := `x * 1 + (2 / 3) > 100`
	e := NewExpressionParser().ParseString(&Query{}, str)
	if e.String() != `((x * 1) + (2 / 3)) > 100` {
		t.Fatal("Unexpected:", "'"+e.String()+"'")
	}
}
