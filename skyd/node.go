package skyd

import (
	"errors"
	"regexp"
)

//------------------------------------------------------------------------------
//
// Globals
//
//------------------------------------------------------------------------------

const NodeIdLength = 7

var ValidNodeIdRegexp *regexp.Regexp = regexp.MustCompile(`^\w+$`)
var ValidNodeHostRegexp *regexp.Regexp = regexp.MustCompile(`^[a-zA-Z0-9\._-]+$`)

var InvalidNodeIdError = errors.New("Invalid node identifier")
var InvalidNodeHostError = errors.New("Invalid node host")
var InvalidNodePortError = errors.New("Invalid node port")

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// A node represents a single server in the cluster.
type Node struct {
	id   string
	host string
	port uint
}

type Nodes []*Node

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// Creates a new node.
func NewNode(id string, host string, port uint) *Node {
	return &Node{
		id:   id,
		host: host,
		port: port,
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
func NewNodeId() string {
	return NewId(NodeIdLength)
}

//--------------------------------------
// Serialization
//--------------------------------------

// Serializes the node into a map that can be converted to JSON.
func (n *Node) Serialize() map[string]interface{} {
	return map[string]interface{}{
		"id":   n.id,
		"host": n.host,
		"port": n.port,
	}
}

//--------------------------------------
// Validation
//--------------------------------------

// Checks that the node is in a valid state.
func (n *Node) Validate() error {
	if !ValidNodeIdRegexp.MatchString(n.id) {
		return InvalidNodeIdError
	}
	if !ValidNodeHostRegexp.MatchString(n.host) {
		return InvalidNodeHostError
	}
	if n.port == 0 {
		return InvalidNodePortError
	}
	return nil
}

//--------------------------------------
// Sorting
//--------------------------------------

func (s Nodes) Len() int {
	return len(s)
}
func (s Nodes) Less(i, j int) bool {
	return s[i].id < s[j].id
}
func (s Nodes) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
