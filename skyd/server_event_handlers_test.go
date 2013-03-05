package skyd

import (
  "testing"
)

// Ensure that we can put an event on the server.
func TestServerPutEvent(t *testing.T) {
  runTestServer(func(s *Server) {
    setupTestTable("foo")
    setupTestProperty("foo", "bar", "object", "string")
    setupTestProperty("foo", "baz", "action", "integer")
    resp, _ := sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/xyz/events", "application/json", `{"timestamp":"2012-01-01T02:00:00Z", "data":{"bar":"myValue"}, "action":{"baz":12}}`)
    assertResponse(t, resp, 200, "", "PUT /tables/:name/objects/:objectId/events failed.")
  })
}
