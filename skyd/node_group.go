package skyd

import (
	"errors"
	"sort"
)

//------------------------------------------------------------------------------
//
// Globals
//
//------------------------------------------------------------------------------

const NodeGroupIdLength = 7

var NodeRequiredError error = errors.New("Node required")
var NodeNotFoundError error = errors.New("Node not found")
var AttachedShardsError error = errors.New("Cannot delete node while group has shards attached")

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A node group clusters organizes nodes inside a cluster. The first node
// in the group acts as the master and replicates to the remaining nodes.
type NodeGroup struct {
	id           string
	leaderNodeId string
	shards       []int
	nodes        []*Node
}

type NodeGroups []*NodeGroup

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// Creates a new node group.
func NewNodeGroup(id string) *NodeGroup {
	return &NodeGroup{
		id:     id,
		shards: make([]int, 0),
		nodes:  make([]*Node, 0),
	}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Identifiers
//--------------------------------------

// Generates a random node identifier.
func NewNodeGroupId() string {
	return NewId(NodeGroupIdLength)
}

//--------------------------------------
// Membership
//--------------------------------------

// Retrieves the leader of the group.
func (g *NodeGroup) leader() *Node {
	if len(g.nodes) > 0 {
		return g.nodes[0]
	}
	return nil
}

// Retrieves the non-leaders in the group.
func (g *NodeGroup) followers() []*Node {
	if len(g.nodes) > 0 {
		return g.nodes[1:]
	}
	return nil
}

//--------------------------------------
// Nodes
//--------------------------------------

// Retrieves a node by id from the group.
func (g *NodeGroup) getNode(id string) *Node {
	for _, node := range g.nodes {
		if node.id == id {
			return node
		}
	}
	return nil
}

// Adds a node to the group.
func (g *NodeGroup) addNode(node *Node) error {
	if node == nil {
		return NodeRequiredError
	}
	if err := node.Validate(); err != nil {
		return err
	}

	g.nodes = append(g.nodes, node)
	sort.Sort(Nodes(g.nodes))

	// Promote to leader if it's the only node.
	if len(g.nodes) == 1 {
		g.leaderNodeId = node.id
	}

	return nil
}

// Removes a node from the group. A node cannot be removed if it is the last
// node in the group and the group still has remaining shards.
func (g *NodeGroup) removeNode(node *Node) error {
	if node == nil {
		return NodeRequiredError
	}

	// Remove from the group.
	for index, n := range g.nodes {
		if n == node {
			// Require that at least one node be in the group if there are shards attached.
			if len(g.nodes) == 1 && len(g.shards) > 0 {
				return AttachedShardsError
			}

			// Otherwise remove the node.
			g.nodes = append(g.nodes[:index], g.nodes[index+1:]...)

			// Promote another node if this node was leader.
			if g.leaderNodeId == node.id {
				if len(g.nodes) == 0 {
					g.leaderNodeId = ""
				} else {
					g.leaderNodeId = g.nodes[0].id
				}
			}
			return nil
		}
	}

	return NodeNotFoundError
}

//--------------------------------------
// Serialization
//--------------------------------------

func (g *NodeGroup) Serialize() map[string]interface{} {
	nodes := []interface{}{}
	for _, node := range g.nodes {
		nodes = append(nodes, node.Serialize())
	}

	shards := make([]int, len(g.shards))
	copy(shards, g.shards)

	return map[string]interface{}{
		"leaderNodeId": g.leaderNodeId,
		"shards":       shards,
		"nodes":        nodes,
	}
}

//--------------------------------------
// Sorting
//--------------------------------------

func (s NodeGroups) Len() int {
	return len(s)
}
func (s NodeGroups) Less(i, j int) bool {
	return s[i].id < s[j].id
}
func (s NodeGroups) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
