package query

import (
	"testing"
)

func TestParser(t *testing.T) {
	str := `SELECT`
	parser := NewParser()
	parser.ParseString(str)
}
