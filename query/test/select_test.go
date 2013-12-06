package test

import (
	"testing"

	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
	"github.com/stretchr/testify/assert"
)

func TestSelectCount(t *testing.T) {
	query := `FOR EACH EVENT SELECT count() END`
	result, err := runDBMapper(query, ast.VarDecls{
		ast.NewVarDecl(1, "foo", "integer"),
	}, map[string][]*core.Event{
		"foo": []*core.Event{
			testevent("2000-01-01T00:00:00Z", 1, 10),
			testevent("2000-01-01T00:00:02Z", 1, 20),
		},
	})
	assert.NoError(t, err)
	assert.NotNil(t, result)
}
