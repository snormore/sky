package parser

import (
	"testing"
)

func TestParserSelectCount(t *testing.T) {
	str := `SELECT count() AS count`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectUnnamedExpression(t *testing.T) {
	str := `SELECT sum(@unit_price)`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectDistinctExpression(t *testing.T) {
	str := `SELECT count(DISTINCT @unit_price)`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectNonAggregatedExpression(t *testing.T) {
	str := `SELECT @name, @price AS unit_price`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectDimensions(t *testing.T) {
	str := `SELECT count() AS count GROUP BY @foo, @bar`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSelectInto(t *testing.T) {
	str := `SELECT count() AS count INTO "xxx\"'"`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserCondition(t *testing.T) {
	str := `WHEN @action == "signup\"'" THEN` + "\n" + `  SELECT count() AS count` + "\n" + `END`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserConditionWithin(t *testing.T) {
	str := `WHEN @action == "signup" WITHIN 1 .. 2 STEPS THEN` + "\n" + `  SELECT count() AS count` + "\n" + `END`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserVariable(t *testing.T) {
	str := `DECLARE @foo AS FLOAT` + "\n"
	str += `DECLARE @bar AS FACTOR(@bat)` + "\n"
	str += `SET @foo = 12` + "\n"
	str += `SELECT sum(@foo) AS total`

	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserSystemVariable(t *testing.T) {
	str := `SELECT count() AS count GROUP BY @@eos`
	query, err := NewParser().ParseString(str)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if query.String() != str {
		t.Fatal("Unexpected:", "'"+query.String()+"'")
	}
}

func TestParserTemporalLoop(t *testing.T) {
	str := `FOR @i EVERY 1 DAY WITHIN 30 DAYS` + "\n"
	str += `  FOR EACH SESSION DELIMITED BY 2 HOURS` + "\n"
	str += `    FOR EACH EVENT` + "\n"
	str += `      SELECT count() AS count` + "\n"
	str += `    END` + "\n"
	str += `  END` + "\n"
	str += `END`

	if query, err := NewParser().ParseString(str); err != nil {
		t.Fatal("Parse error:", err)
	} else if query.String() != str {
		t.Fatal("\nExpected:\n" + str + "\n\nGot:\n" + query.String())
	}
}

func TestParserError(t *testing.T) {
	str := `SELECT count() AS count` + "\n"
	str += `EXIT`

	if query, err := NewParser().ParseString(str); err != nil {
		t.Fatal("Parse error:", err)
	} else if query.String() != str {
		t.Fatal("\nExpected:\n" + str + "\n\nGot:\n" + query.String())
	}
}

func TestParserDebug(t *testing.T) {
	str := `DEBUG(@x + 12)`
	if query, err := NewParser().ParseString(str); err != nil {
		t.Fatal("Parse error:", err)
	} else if query.String() != str {
		t.Fatal("\nExpected:\n" + str + "\n\nGot:\n" + query.String())
	}
}

func TestParserSyntaxError(t *testing.T) {
	_, err := NewParser().ParseString(`SELECT count() AS` + "\n" + `count GRP BY @action`)
	if err == nil || err.Error() != "Unexpected 'GRP' at line 2, char 8, syntax error" {
		t.Fatal("Unexpected parse error: ", err)
	}
}
