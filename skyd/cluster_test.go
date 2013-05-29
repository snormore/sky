package skyd

import (
	"math/rand"
	"testing"
)

// Ensure that node groups can be added to a cluster.
func TestClusterAddNodeGroups(t *testing.T) {
	rand.Seed(1)
	cluster := NewCluster()
	g0 := NewNodeGroup(NewNodeGroupId())
	cluster.AddNodeGroup(g0)
	g1 := NewNodeGroup(NewNodeGroupId())
	cluster.AddNodeGroup(g1)

	if len(cluster.groups) != 2 {
		t.Fatalf("Unexpected group count: %v", len(cluster.groups))
	}
}

// Ensure that nodes can be added to a group.
func TestClusterAddNodesToGroups(t *testing.T) {
	rand.Seed(2)
	cluster := NewCluster()
	g0 := NewNodeGroup(NewNodeGroupId())
	cluster.AddNodeGroup(g0)
	n0 := NewNode(NewNodeId(), "127.0.0.1", 8585)
	n1 := NewNode(NewNodeId(), "localhost", 8586)

	if err := cluster.AddNode(n0, g0); err != nil {
		t.Fatalf("Unable to add node to group: %v", err)
	}
	if err := cluster.AddNode(n1, g0); err != nil {
		t.Fatalf("Unable to add node to same group: %v", err)
	}
	if g0.leaderNodeId != n0.id {
		t.Fatalf("Unexpected group leader: %v", g0.leaderNodeId)
	}
}

// Ensure that nodes can be added to a group.
func TestClusterRemoveNodesFromGroup(t *testing.T) {
	rand.Seed(3)
	cluster := NewCluster()
	g0 := NewNodeGroup(NewNodeGroupId())
	cluster.AddNodeGroup(g0)
	n0 := NewNode(NewNodeId(), "127.0.0.1", 8585)
	cluster.AddNode(n0, g0)
	n1 := NewNode(NewNodeId(), "localhost", 8586)
	cluster.AddNode(n1, g0)

	if err := cluster.RemoveNode(n0); err != nil {
		t.Fatalf("Unable to remove node: %v", err)
	}
}
