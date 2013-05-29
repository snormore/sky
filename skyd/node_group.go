package skyd

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A node group clusters organizes nodes inside a cluster. The first node
// in the group acts as the master and replicates to the remaining nodes.
type NodeGroup struct {
	shards []uint8
	nodes  []*Node
}

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// Creates a new node group.
func NewNodeGroup() *NodeGroup {
	return &NodeGroup{}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Membership
//--------------------------------------

// Retrieves the leader of the group.
func (g *NodeGroup) Leader() *Node {
	if len(g.nodes) > 0 {
		return g.nodes[0]
	}
	return nil
}

// Retrieves the non-leaders in the group.
func (g *NodeGroup) Followers() []*Node {
	if len(g.nodes) > 0 {
		return g.nodes[1:]
	}
	return nil
}

//--------------------------------------
// Nodes
//--------------------------------------

// Adds a node to the group.
func (g *NodeGroup) addNode(name string) error {
	if name == "" {
		return InvalidNodeNameError
	}

	g.nodes = append(g.nodes, NewNode(name))
	return nil
}

// Checks if the group contains a node with a given name.
func (g *NodeGroup) hasNode(name string) bool {
	for _, node := range g.nodes {
		if node.name == name {
			return true
		}
	}
	return false
}

//--------------------------------------
// Serialization
//--------------------------------------

func (g *NodeGroup) Serialize() map[string]interface{} {
	nodes := []interface{}{}
	for _, node := range g.nodes {
		nodes = append(nodes, node.Serialize())
	}

	shards := make([]uint8, len(g.shards))
	copy(shards, g.shards)
	return map[string]interface{}{
		"shards": shards,
		"nodes":  nodes,
	}
}
