package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure that we can retrieve a list of all available tables on the server.
func TestServerTableGetAll(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestTable("bar")
		code, resp := getJSON("/tables")
		assert.Equal(t, code, 200)
		if assert.NotNil(t, resp) {
			resp := resp.([]interface{})
			assert.Equal(t, len(resp), 2)
			assert.Equal(t, resp[0].(map[string]interface{})["name"], "bar")
			assert.Equal(t, resp[1].(map[string]interface{})["name"], "foo")
		}
	})
}

// Ensure that we can retrieve a single table on the server.
func TestServerTableGet(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		code, resp := getJSON("/tables/foo")
		assert.Equal(t, code, 200)
		if assert.NotNil(t, resp) {
			resp := resp.(map[string]interface{})
			assert.Equal(t, resp["name"], "foo")
		}
	})
}

// Ensure that we can create a new table through the server.
func TestServerTableCreate(t *testing.T) {
	runTestServer(func(s *Server) {
		code, resp := postJSON("/tables", `{"name":"foo"}`)
		assert.Equal(t, code, 200)
		if assert.NotNil(t, resp) {
			resp := resp.(map[string]interface{})
			assert.Equal(t, resp["name"], "foo")
		}
	})
}

// Ensure that we can delete a table through the server.
func TestServerTableDelete(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		code, resp := deleteJSON("/tables/foo", "")
		assert.Equal(t, code, 200)
		assert.Nil(t, resp)
	})
}

// Ensure that we can retrieve a list of object keys from the server for a table.
func TestServerTableKeys(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "value", true, "integer")
		setupTestData(t, "foo", [][]string{
			[]string{"a0", "2012-01-01T00:00:00Z", `{"data":{"value":1}}`},
			[]string{"a1", "2012-01-01T00:00:00Z", `{"data":{"value":2}}`},
			[]string{"a1", "2012-01-01T00:00:01Z", `{"data":{"value":3}}`},
			[]string{"a2", "2012-01-01T00:00:00Z", `{"data":{"value":4}}`},
			[]string{"a2", "2012-01-01T00:00:01Z", `{"data":{"value":4}}`},
			[]string{"a3", "2012-01-01T00:00:00Z", `{"data":{"value":5}}`},
		})

		setupTestTable("bar")
		setupTestProperty("bar", "value", true, "integer")
		setupTestData(t, "bar", [][]string{
			[]string{"b0", "2012-01-01T00:00:00Z", `{"data":{"value":1}}`},
			[]string{"b1", "2012-01-01T00:00:00Z", `{"data":{"value":2}}`},
			[]string{"b1", "2012-01-01T00:00:01Z", `{"data":{"value":3}}`},
			[]string{"b2", "2012-01-01T00:00:00Z", `{"data":{"value":4}}`},
			[]string{"b3", "2012-01-01T00:00:00Z", `{"data":{"value":5}}`},
		})

		code, resp := getJSON("/tables/foo/keys")
		assert.Equal(t, code, 200)
		if assert.NotNil(t, resp) {
			resp := resp.([]interface{})
			if assert.Equal(t, len(resp), 4) {
				assert.Equal(t, resp[0], "a0")
				assert.Equal(t, resp[1], "a1")
				assert.Equal(t, resp[2], "a2")
				assert.Equal(t, resp[3], "a3")
			}
		}
	})
}
