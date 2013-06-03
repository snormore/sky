package skyd

import (
	"testing"
	"time"
)

// Ensure that we can add a server to the cluster and the configuration is
// replicated between the nodes.
func TestMultinodeJoin(t *testing.T) {
	f0 := func(s *Server) {
		warn("1? [%p] %d", s, s.port)
		time.Sleep(100 * time.Millisecond)
		if len(s.cluster.groups) != 1 {
			t.Fatalf("Unexpected group count: %v", len(s.cluster.groups))
		}
		if len(s.cluster.groups[0].nodes) != 2 {
			t.Fatalf("Unexpected node count: [%p] %v", s, s.cluster.groups[0].nodes)
		}
	}
	f1 := func(s *Server) {
		warn("2? [%p] %d", s, s.port)
		if err := s.Join("localhost", 8800); err != nil {
			t.Fatalf("Unable to join: %v", err)
		}
	}
	runTestServers(f0, f1)
}
