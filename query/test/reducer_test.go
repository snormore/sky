package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/skydb/sky/core"
	"github.com/skydb/sky/db"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hashmap"
	"github.com/skydb/sky/query/codegen/mapper"
	"github.com/skydb/sky/query/parser"
	"github.com/stretchr/testify/assert"
)

func TestReducerSelectCount(t *testing.T) {
	query := `
		FOR EACH EVENT
			SELECT count()
		END
	`
	result, err := runDBMapper(query, ast.VarDecls{
		ast.NewVarDecl(1, "foo", "integer"),
	}, map[string][]*core.Event{
		"foo": []*core.Event{
			testevent("2000-01-01T00:00:00Z", 1, 10),
			testevent("2000-01-01T00:00:02Z", 1, 20),
		},
		"bar": []*core.Event{
			testevent("2000-01-01T00:00:00Z", 1, 40),
		},
	})
	assert.NoError(t, err)
	if assert.NotNil(t, result) {
		assert.Equal(t, result.Get(HASH_COUNT), 3)
	}
}

// Executes a query against a given set of data, reduces it and return the results.
func runDBMapReducer(query string, decls ast.VarDecls, objects map[string][]*core.Event) (map[string]interface{}, error) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)

	db := db.New(path, 1, false, 4096, 126)
	if err := db.Open(); err != nil {
		return nil, err
	}
	defer db.Close()

	// Insert into db.
	if _, err := db.InsertObjects("TBL", objects); err != nil {
		return nil, err
	}

	// Retrieve cursors.
	cursors, err := db.Cursors("TBL")
	if err != nil {
		return nil, err
	}
	defer cursors.Close()

	// Create a query.
	q, err := parser.ParseString(query)
	if err != nil {
		return nil, err
	}
	q.DeclaredVarDecls = append(q.DeclaredVarDecls, decls...)

	// Setup factor test data.
	f := db.TableFactorizer("TBL")
	f.Factorize("action", "A0", true)
	f.Factorize("action", "A1", true)

	// Create a mapper generated from the query.
	m, err := mapper.New(q, f)
	if err != nil {
		return nil, err
	}
	// m.Dump()

	// Execute the mapper.
	result := hashmap.New()
	if err = m.Execute(cursors[0], "", result); err != nil {
		return nil, err
	}

	return result, nil
}
