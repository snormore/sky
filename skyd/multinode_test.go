package skyd

import (
	"testing"
)

// Ensure that we can add a server to the cluster and the configuration is
// replicated between the nodes.
func TestMultinodeJoin(t *testing.T) {
	runTestServers(func(s *Server) {
		resp, err := sendTestHttpRequest("GET", "http://localhost:8800/ping", "application/json", "")
		if err != nil {
			t.Fatalf("Unable to ping: %v", err)
		}
		assertResponse(t, resp, 200, `{"message":"ok"}`+"\n", "GET /ping failed.")
	},
		func(s *Server) {
			resp, err := sendTestHttpRequest("GET", "http://localhost:8801/ping", "application/json", "")
			if err != nil {
				t.Fatalf("Unable to ping: %v", err)
			}
			assertResponse(t, resp, 200, `{"message":"ok"}`+"\n", "GET /ping failed.")
		})
}
