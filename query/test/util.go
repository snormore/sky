package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/skydb/sky/core"
	"github.com/skydb/sky/db"
	"github.com/skydb/sky/query/codegen/mapper"
	"github.com/skydb/sky/query/parser"
)

// Executes a query against a given set of data and return the results.
func runDBMapper(query string, objects map[string][]*core.Event) (interface{}, error) {
	path, _ := ioutil.TempDir("", "")
	defer os.RemoveAll(path)

	db := db.New(path, 1, false, 4096, 126)
	if err := db.Open(); err != nil {
		debugln("run.mapper.!")
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

	// Create a mapper generated from the query.
	m, err := mapper.New(q)
	if err != nil {
		return nil, err
	}
	m.Dump()

	// Execute the mapper.
	return m.Execute(cursors[0], "", nil), nil
}

func testevent(timestamp string, args ...interface{}) *core.Event {
	e := &core.Event{Timestamp: musttime(timestamp)}
	e.Data = make(map[int64]interface{})
	for i := 0; i < len(args); i += 2 {
		key := args[i].(int)
		e.Data[int64(key)] = args[i+1]
	}
	return e
}

func musttime(timestamp string) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		panic(err)
	}
	return t
}

func debugln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
