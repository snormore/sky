package skyd

import (
	"fmt"
	"os"
	"testing"
)

// Ensure that we can create a new table through the server.
func TestServerCreateTable(t *testing.T) {
	runTestServer(func(s *Server) {
		resp, err := sendTestHttpRequest("POST", "http://localhost:8585/tables", "application/json", `{"name":"foo"}`)
		if err != nil {
			t.Fatalf("Unable to create table: %v", err)
		}
		assertResponse(t, resp, 200, "", "POST /tables failed.")
		if _, err := os.Stat(fmt.Sprintf("%v/tables/foo", s.Path())); os.IsNotExist(err) {
			t.Fatalf("POST /tables did not create table.")
		}
	})
}

// Ensure that we can delete a table through the server.
func TestServerDeleteTable(t *testing.T) {
	runTestServer(func(s *Server) {
		// Create table.
		resp, err := sendTestHttpRequest("POST", "http://localhost:8585/tables", "application/json", `{"name":"foo"}`)
		if err != nil {
			t.Fatalf("Unable to create table: %v", err)
		}
		assertResponse(t, resp, 200, "", "POST /tables failed.")
		if _, err := os.Stat(fmt.Sprintf("%v/tables/foo", s.Path())); os.IsNotExist(err) {
			t.Fatalf("POST /tables did not create table.")
		}

		// Delete table.
		resp, _ = sendTestHttpRequest("DELETE", "http://localhost:8585/tables/foo", "application/json", ``)
		assertResponse(t, resp, 200, "", "DELETE /tables/:name failed.")
		if _, err := os.Stat(fmt.Sprintf("%v/tables/foo", s.Path())); !os.IsNotExist(err) {
			t.Fatalf("DELETE /tables/:name did not delete table.")
		}
	})
}
