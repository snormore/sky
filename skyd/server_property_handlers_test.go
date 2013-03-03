package skyd

import (
	"testing"
)

// Ensure that we can create a property through the server.
func TestServerCreateProperty(t *testing.T) {
  runTestServer(func() {
    setupTestTable("foo")
    resp, err := sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/properties", "application/json", `{"name":"bar", "type":"object", "dataType":"string"}`)
    if err != nil {
  		t.Fatalf("Unable to create property: %v", err)
    }
    assertResponse(t, resp, 200, `{"id":1,"name":"bar","type":"object","dataType":"string"}`+"\n", "POST /tables/:name/properties failed.")
  })
}

// Ensure that we can retrieve all properties through the server.
func TestServerGetProperties(t *testing.T) {
  runTestServer(func() {
    setupTestTable("foo")
    setupTestProperty("foo", "bar", "object", "string")
    setupTestProperty("foo", "baz", "action", "integer")
    resp, err := sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/properties", "application/json", "")
    if err != nil {
  		t.Fatalf("Unable to get properties: %v", err)
    }
    assertResponse(t, resp, 200, `[{"id":-1,"name":"baz","type":"action","dataType":"integer"},{"id":1,"name":"bar","type":"object","dataType":"string"}]`+"\n", "GET /tables/:name/properties failed.")
  })
}

