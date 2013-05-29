package skyd

import (
	"math/rand"
	"testing"
)

//------------------------------------------------------------------------------
//
// Tests
//
//------------------------------------------------------------------------------

//--------------------------------------
// Node Groups
//--------------------------------------

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
	if len(g0.shards) != 256 {
		t.Fatalf("Unexpected shard count: %v", len(g0.shards))
	}
	if len(g1.shards) != 0 {
		t.Fatalf("Unexpected shard count for secondary group: %v", len(g1.shards))
	}
}

// Ensure that node groups can be removed from a cluster.
func TestClusterRemoveNodeGroups(t *testing.T) {
	cluster := NewCluster()
	g0 := NewNodeGroup("a")
	cluster.AddNodeGroup(g0)
	g1 := NewNodeGroup("b")
	cluster.AddNodeGroup(g1)

	if err := cluster.RemoveNodeGroup(g1); err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if err := cluster.RemoveNodeGroup(g0); err != NodeGroupAttachedShardsError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestClusterAddMissingNodeGroup(t *testing.T) {
	if err := NewCluster().AddNodeGroup(nil); err != NodeGroupRequiredError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestClusterAddNamelessNodeGroup(t *testing.T) {
	if err := NewCluster().AddNodeGroup(NewNodeGroup("")); err != InvalidNodeGroupIdError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestClusterAddGroupWithExistingShards(t *testing.T) {
	group := NewNodeGroup(NewNodeGroupId())
	group.shards = append(group.shards, 0)
	if err := NewCluster().AddNodeGroup(group); err != ExistingShardsError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestClusterRemoveMissingNodeGroup(t *testing.T) {
	if err := NewCluster().RemoveNodeGroup(nil); err != NodeGroupRequiredError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestClusterRemoveNodeGroupWithNodes(t *testing.T) {
	cluster := NewCluster()
	g0 := NewNodeGroup("a")
	cluster.AddNodeGroup(g0)
	n0 := NewNode(NewNodeId(), "localhost", 8585)
	cluster.AddNode(n0, g0)

	if err := cluster.RemoveNodeGroup(g0); err != NodeGroupAttachedNodesError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestClusterRemoveNonExistentNodeGroup(t *testing.T) {
	if err := NewCluster().RemoveNodeGroup(NewNodeGroup("a")); err != NodeGroupNotFoundError {
		t.Fatalf("Unexpected error: %v", err)
	}
}

//--------------------------------------
// Nodes
//--------------------------------------

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
