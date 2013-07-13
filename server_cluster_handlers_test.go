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
		warn("0.1")
		time.Sleep(TestHeartbeatTimeout)
		warn("0.2")
		if num := len(s.cluster.groups[0].nodes); num != 1 {
			t.Fatalf("[%d.%p] Unexpected node count: %v", 0, s, num)
		}
		warn("0.3")
		if num := s.clusterRaftServer.MemberCount(); num != 1 {
			t.Fatalf("[%d.%p] Unexpected cluster member count: %v", 0, s, num)
		}
		warn("0.4")
	}
	f1 := func(s *Server) {
		warn("1.1")
		if err := s.Leave(); err != nil {
			t.Fatalf("Unable to leave cluster: %v", err)
		}
		warn("1.2")
		time.Sleep(TestHeartbeatTimeout)
		warn("1.3")
		if s.Running() {
			t.Fatalf("[%d.%p] Unexpected server state: running", 1, s)
		}
		warn("1.4")
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
		err := rpc("localhost", 8800, "POST", "/cluster/groups", map[string]interface{}{"nodeGroupId": "foo"}, resp)
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

// Ensure that we can't remove a group that has nodes.
func TestHttpClusterRemoveNodeGroupWithNodes(t *testing.T) {
	f0 := func(s *Server) {
		nodeGroupId := s.cluster.groups[0].id
		err := rpc("localhost", 8800, "DELETE", "/cluster/groups/"+nodeGroupId, nil, nil)
		if err == nil || strings.Index(err.Error(), "Cannot delete node group while nodes are attached") != 0 {
			t.Fatalf("Unexpected error: %v", err)
		}
		if num := len(s.cluster.groups); num != 1 {
			t.Fatalf("[%p] Unexpected cluster state: %d", s, num)
		}
	}
	runTestServers(true, f0)
}
