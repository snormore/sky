package mapper

import (
	"testing"

	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

func TestMapper(t *testing.T) {
	m, err := New(ast.NewQuery())
	assert.NoError(t, err)
	m.Dump()
	assert.Equal(t, 12, m.Execute())
}
