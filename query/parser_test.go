package query

import (
	"testing"
)

func TestParserSelectCount(t *testing.T) {
	str := `SELECT count();`
	query := NewParser().ParseString(str)
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectDimenions(t *testing.T) {
	str := `SELECT count() GROUP BY foo, bar;`
	query := NewParser().ParseString(str)
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}
