package skyd

import (
	"errors"
	"sync"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// The cluster manages the topology of servers for a distributed Sky database.
// Clusters are made up of cluster groups which are sets of servers that
// manage a subset of the total dataset.
type Cluster struct {
	groups           []*NodeGroup
	mutex            sync.Mutex
}

var NodeGroupNotFoundError = errors.New("Node group not found")
var DuplicateNodeError = errors.New("Duplicate node already exists")

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

// Creates a new cluster.
func NewCluster() *Cluster {
	return &Cluster{
		groups: []*NodeGroup{},
	}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Node Groups
//--------------------------------------

// Adds a group to the cluster.
func (c *Cluster) addNodeGroup() {
	c.mutex.Lock()
	c.mutex.Unlock()
	c.groups = append(c.groups, NewNodeGroup())
}

// Removes a group from the cluster.
func (c *Cluster) removeNodeGroup(index int) error {
	c.mutex.Lock()
	c.mutex.Unlock()
	
	if index < 0 || len(c.groups) <= index {
		return NodeGroupNotFoundError
	}
	
	c.groups = append(c.groups[:index], c.groups[index+1:]...)
	return nil
}

//--------------------------------------
// Nodes
//--------------------------------------

// Adds a node to an existing group in the cluster.
func (c *Cluster) addNode(name string, groupIndex int) error {
	c.mutex.Lock()
	c.mutex.Unlock()

	// Validate that the group exists.
	if groupIndex < 0 || len(c.groups) <= groupIndex {
		return NodeGroupNotFoundError
	}
	
	// Check if the node exists in the cluster already.
	for _, group := range c.groups {
		if group.hasNode(name) {
			return DuplicateNodeError
		}
	}

	group := c.groups[groupIndex]
	return group.addNode(name)
}

//--------------------------------------
// Serialization
//--------------------------------------

// Converts the cluster topology to an object that can be easily serialized
// to JSON outside the cluster lock.
func (c *Cluster) serialize() map[string]interface{} {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Serialize groups.
	groups := []interface{}{}
	for _, group := range c.groups {
		groups = append(groups, group.Serialize())
	}
	
	return map[string]interface{}{"groups":groups}
}

