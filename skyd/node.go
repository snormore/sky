package skyd

import (
	"errors"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A node represents a single server in the cluster.
type Node struct {
	name string
}

var InvalidNodeNameError = errors.New("Invalid node name")

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// Creates a new node.
func NewNode(name string) *Node {
	return &Node{name: name}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Serialization
//--------------------------------------

func (n *Node) Serialize() map[string]interface{} {
	return map[string]interface{}{"name": n.name}
}
