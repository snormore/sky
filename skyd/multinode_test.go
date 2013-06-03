package skyd

import (
	"testing"
	"time"
)

// Ensure that we can add a server to the cluster and the configuration is
// replicated between the nodes.
func TestMultinodeJoin(t *testing.T) {
	f0 := func(s *Server) {
		// Wait for #2 to join
		time.Sleep(100 * time.Millisecond)
		
		// Make sure that they're added to the group & to the cluster peers.
		if num := len(s.cluster.groups[0].nodes); num != 2 {
			t.Fatalf("[1.%p] Unexpected node count: %v", s, num)
		}
		if num := s.clusterRaftServer.MemberCount(); num != 2 {
			t.Fatalf("[1.%p] Unexpected cluster member count: %v", s, num)
		}
	}
	f1 := func(s *Server) {
		if err := s.Join("localhost", 8800); err != nil {
			t.Fatalf("Unable to join: %v", err)
		}
		// Wait for 
		time.Sleep(50 * time.Millisecond)
		
		// Make sure that cluster's log is replicated back to #2.
		if num := len(s.cluster.groups[0].nodes); num != 2 {
			t.Fatalf("[2.%p] Unexpected node count: %v", s, num)
		}
		if num := s.clusterRaftServer.MemberCount(); num != 2 {
			t.Fatalf("[2.%p] Unexpected cluster member count: %v", s, num)
		}
	}
	runTestServers(f0, f1)
}
