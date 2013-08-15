package query

import (
	"testing"
)

func TestParser(t *testing.T) {
	str := `SELECT xyz;`
	parser := NewParser()
	parser.ParseString(str)
}
