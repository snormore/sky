package skyd

import (
	"testing"
	"time"
)

// Ensure that we can add a server to the cluster and the configuration is
// replicated between the nodes.
func TestMultinodeJoin(t *testing.T) {
	f0 := func(s *Server) {
		time.Sleep(100 * time.Millisecond)
		if len(s.cluster.groups[0].nodes) != 2 {
			t.Fatalf("Unexpected node count: %v", len(s.cluster.groups[0].nodes))
		}
	}
	f1 := func(s *Server) {
		resp, err := sendTestHttpRequest("POST", "http://localhost:8800/cluster/nodes", "application/json", `{"host":"localhost","port":8801}`)
		if err != nil {
			t.Fatalf("Unable to join: %v", err)
		}
		assertResponse(t, resp, 200, ``+"\n", "POST /cluster/nodes failed.")
	}
	runTestServers(f0, f1)
}
