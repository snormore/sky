package skyd

import (
	"strings"
	"testing"
	"time"
)

//------------------------------------------------------------------------------
//
// Tests
//
//------------------------------------------------------------------------------

//--------------------------------------
// Nodes
//--------------------------------------

// Ensure that we can add a server to the cluster and the configuration is
// replicated between the nodes.
func TestHttpClusterAddNode(t *testing.T) {
	assert := func(index int, s *Server) {
		// Make sure that they're added to the group & to the cluster peers.
		if num := len(s.cluster.groups[0].nodes); num != 2 {
			t.Fatalf("[%d.%p] Unexpected node count: %v", index, s, num)
		}
		if num := s.clusterRaftServer.MemberCount(); num != 2 {
			t.Fatalf("[%d.%p] Unexpected cluster member count: %v", index, s, num)
		}
	}
	
	f0 := func(s *Server) {
		// Wait for #2 to join
		time.Sleep(100 * time.Millisecond)
		assert(0, s)
	}
	f1 := func(s *Server) {
		if err := s.Join("localhost", 8800); err != nil {
			t.Fatalf("Unable to join cluster: %v", err)
		}
	
		time.Sleep(100 * time.Millisecond)
		assert(1, s)
	}
	runTestServers(false, f0, f1)
}

// Ensure that we can remove a server from the cluster.
func TestHttpClusterRemoveNode(t *testing.T) {
	f0 := func(s *Server) {
		time.Sleep(TestHeartbeatTimeout)
		if num := len(s.cluster.groups[0].nodes); num != 1 {
			t.Fatalf("[%d.%p] Unexpected node count: %v", 0, s, num)
		}
		if num := s.clusterRaftServer.MemberCount(); num != 1 {
			t.Fatalf("[%d.%p] Unexpected cluster member count: %v", 0, s, num)
		}
	}
	f1 := func(s *Server) {
		if err := s.Leave(); err != nil {
			t.Fatalf("Unable to leave cluster: %v", err)
		}
		time.Sleep(TestHeartbeatTimeout)
		if s.Running() {
			t.Fatalf("[%d.%p] Unexpected server state: running", 1, s)
		}
	}
	runTestServers(true, f0, f1)
}

// Ensure that we cannot remove the last node from a group.
func TestHttpClusterRemoveLastNode(t *testing.T) {
	f0 := func(s *Server) {
		if err := s.Leave(); strings.Index(err.Error(), "Cannot remove last node from group") != 0 {
			t.Fatalf("Last node should not be able to leave cluster", err)
		}
		if !s.Running() {
			t.Fatalf("[%p] Unexpected server state: stopped", s)
		}
	}
	runTestServers(true, f0)
}


//--------------------------------------
// Groups
//--------------------------------------

// Ensure that we can create and remove a group.
func TestHttpClusterCreateAndRemoveNodeGroup(t *testing.T) {
	f0 := func(s *Server) {
		// Add group.
		resp := make(map[string]interface{})
		err := rpc("localhost", 8800, "POST", "/cluster/groups", map[string]interface{}{"nodeGroupId":"foo"}, resp)
		if err != nil {
			t.Fatalf("Unable to add group: %v", err)
		}
		if num := len(s.cluster.groups); num != 2 || s.cluster.groups[1].id != "foo" {
			t.Fatalf("[%p] Unexpected group state: %d", s, num)
		}

		// Remove group.
		err = rpc("localhost", 8800, "DELETE", "/cluster/groups/foo", nil, nil)
		if err != nil {
			t.Fatalf("Unable to remove group: %v", err)
		}
		if num := len(s.cluster.groups); num != 1 {
			t.Fatalf("[%p] Unexpected group state after removal: %d", s, num)
		}
	}
	runTestServers(true, f0)
}
