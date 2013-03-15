package skyd

import (
	"testing"
)

// Ensure that we can create a property through the server.
func TestServerCreateProperty(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		resp, _ := sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/properties", "application/json", `{"name":"bar", "type":"object", "dataType":"string"}`)
		assertResponse(t, resp, 200, `{"id":1,"name":"bar","type":"object","dataType":"string"}`+"\n", "POST /tables/:name/properties failed.")
	})
}

// Ensure that we can retrieve all properties through the server.
func TestServerGetProperties(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "bar", "object", "string")
		setupTestProperty("foo", "baz", "action", "integer")
		resp, _ := sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/properties", "application/json", "")
		assertResponse(t, resp, 200, `[{"id":-1,"name":"baz","type":"action","dataType":"integer"},{"id":1,"name":"bar","type":"object","dataType":"string"}]`+"\n", "GET /tables/:name/properties failed.")
	})
}

// Ensure that we can retrieve a single property through the server.
func TestServerGetProperty(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "bar", "object", "string")
		setupTestProperty("foo", "baz", "action", "integer")
		resp, _ := sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/properties/bar", "application/json", "")
		assertResponse(t, resp, 200, `{"id":1,"name":"bar","type":"object","dataType":"string"}`+"\n", "GET /tables/:name/properties/:propertyName failed.")
	})
}

// Ensure that we can update a property name through the server.
func TestServerUpdateProperty(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "bar", "object", "string")
		setupTestProperty("foo", "baz", "action", "integer")
		resp, _ := sendTestHttpRequest("PATCH", "http://localhost:8585/tables/foo/properties/bar", "application/json", `{"name":"bat"}`)
		assertResponse(t, resp, 200, `{"id":1,"name":"bat","type":"object","dataType":"string"}`+"\n", "PATCH /tables/:name/properties/:propertyName failed.")
	})
}

// Ensure that we can delete a property on the server.
func TestServerDeleteProperty(t *testing.T) {
	runTestServer(func(s *Server) {
		setupTestTable("foo")
		setupTestProperty("foo", "bar", "object", "string")
		setupTestProperty("foo", "baz", "action", "integer")
		resp, _ := sendTestHttpRequest("DELETE", "http://localhost:8585/tables/foo/properties/bar", "application/json", "")
		assertResponse(t, resp, 200, "", "DELETE /tables/:name/properties/:propertyName failed.")
		resp, _ = sendTestHttpRequest("GET", "http://localhost:8585/tables/foo/properties", "application/json", "")
		assertResponse(t, resp, 200, `[{"id":-1,"name":"baz","type":"action","dataType":"integer"}]`+"\n", "GET /tables/:name/properties after delete failed.")
	})
}
