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

func TestParserCondition(t *testing.T) {
	str := `WHEN action == "signup" THEN` + "\n" + `  SELECT count();` + "\n" + `END`
	query := NewParser().ParseString(str)
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserConditionWithin(t *testing.T) {
	str := `WHEN action == "signup" WITHIN 1 .. 2 STEPS THEN` + "\n" + `  SELECT count();` + "\n" + `END`
	query := NewParser().ParseString(str)
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}
