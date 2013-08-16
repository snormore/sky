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

func TestParserSelectDimensions(t *testing.T) {
	str := `SELECT count() GROUP BY foo, bar;`
	query := NewParser().ParseString(str)
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectInto(t *testing.T) {
	str := `SELECT count() INTO "xxx";`
	query := NewParser().ParseString(str)
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}
