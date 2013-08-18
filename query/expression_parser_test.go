package query

import (
	"testing"
)

func TestExpressionParserSimple(t *testing.T) {
    str := `timestamp >= 2 && timestamp < 6`
    e := NewExpressionParser().ParseString(&Query{}, str)
    if e.String() != `(timestamp >= 2) && (timestamp < 6)` {
        t.Fatal("Unexpected:", "'"+e.String()+"'")
    }
}

func TestExpressionParserComplex(t *testing.T) {
    str := `x * 1 + (2 / 3) > 100`
    e := NewExpressionParser().ParseString(&Query{}, str)
    if e.String() != `((x * 1) + (2 / 3)) > 100` {
        t.Fatal("Unexpected:", "'"+e.String()+"'")
    }
}

