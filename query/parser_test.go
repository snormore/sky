package query

import (
	"testing"
)

func TestParserSelectCount(t *testing.T) {
	str := `SELECT count() AS count;`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectDimensions(t *testing.T) {
	str := `SELECT count() AS count GROUP BY foo, bar;`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectInto(t *testing.T) {
	str := `SELECT count() AS count INTO "xxx";`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserCondition(t *testing.T) {
	str := `WHEN action == "signup" THEN` + "\n" + `  SELECT count() AS count;` + "\n" + `END`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserConditionWithin(t *testing.T) {
	str := `WHEN action == "signup" WITHIN 1 .. 2 STEPS THEN` + "\n" + `  SELECT count() AS count;` + "\n" + `END`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserVariable(t *testing.T) {
	str := `DECLARE foo AS FLOAT` + "\n"
	str += `SET foo = 12` + "\n"
	str += `SELECT sum(foo) AS total;`

	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserError(t *testing.T) {
	_, err := NewParser().ParseString(`SELECT count() AS` + "\n" + `count GRP BY action;`)
	if err == nil || err.Error() != "Unexpected 'GRP' at line 2, char 8, syntax error" {
		t.Fatal("Unexpected parse error: ", err)
	}
}
