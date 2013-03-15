package skyd

import (
  "testing"
)

// Ensure that we can query the server.
func TestServerQuery(t *testing.T) {
  runTestServer(func(s *Server) {
    setupTestTable("foo")
    setupTestProperty("foo", "bar", "object", "string")
    setupTestProperty("foo", "baz", "action", "integer")
    
    // Send events.
    resp, _ := sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/0/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"val1"}}`)
    resp, _ = sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/1/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"val2"}}`)
    resp, _ = sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/2/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"val3"}}`)
    resp, _ = sendTestHttpRequest("PUT", "http://localhost:8585/tables/foo/objects/3/events/2012-01-01T02:00:00Z", "application/json", `{"data":{"bar":"val4"}}`)

    // Run query.
    resp, _ = sendTestHttpRequest("POST", "http://localhost:8585/tables/foo/query", "application/json", ``)
    assertResponse(t, resp, 200, `{"event_count":4,"path_count":4}`+"\n", "POST /tables/:name/query failed.")
  })
}
