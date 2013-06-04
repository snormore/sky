package skyd

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

//------------------------------------------------------------------------------
//
// Globals
//
//------------------------------------------------------------------------------

const (
	ShardCount = 256
)

var NodeGroupRequiredError = errors.New("Node group required")
var NodeGroupNotFoundError = errors.New("Node group not found")
var NodeGroupAttachedShardsError = errors.New("Cannot delete node group while shards are attached")
var NodeGroupAttachedNodesError = errors.New("Cannot delete node group while nodes are attached")
var ExistingShardsError = errors.New("Node group cannot be added with existing shards")
var ShardNotFoundError = errors.New("Shard not found")

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// The cluster manages the topology of servers for a distributed Sky database.
// Clusters are made up of cluster groups which are sets of servers that
// manage a subset of the total dataset.
type Cluster struct {
	groups []*NodeGroup
	mutex  sync.Mutex
}

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

// Finds a group in the cluster by id.
func (c *Cluster) GetNodeGroup(id string) *NodeGroup {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getNodeGroup(id)
}

func (c *Cluster) getNodeGroup(id string) *NodeGroup {
	for _, group := range c.groups {
		if group.id == id {
			return group
		}
	}
	return nil
}

// Finds the most appropriate group to add a new node to.
func (c *Cluster) GetAvailableNodeGroup() *NodeGroup {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Find the least crowded group.
	var group *NodeGroup
	for _, g := range c.groups {
		if group == nil || len(g.nodes) < len(group.nodes) {
			group = g
		}
	}

	return group
}

// Adds a group to the cluster.
func (c *Cluster) AddNodeGroup(group *NodeGroup) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if group == nil {
		return NodeGroupRequiredError
	}
	if err := group.Validate(); err != nil {
		return err
	}
	if len(group.shards) > 0 {
		return ExistingShardsError
	}

	// Append the group.
	c.groups = append(c.groups, group)
	sort.Sort(NodeGroups(c.groups))

	// Add all the shards to the first group.
	if len(c.groups) == 1 {
		for i := 0; i < ShardCount; i++ {
			group.shards = append(group.shards, i)
		}
	}
	return nil
}

// Removes a group from the cluster.
func (c *Cluster) RemoveNodeGroup(group *NodeGroup) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if group == nil {
		return NodeGroupRequiredError
	}
	if len(group.nodes) > 0 {
		return NodeGroupAttachedNodesError
	}
	if len(group.shards) > 0 {
		return NodeGroupAttachedShardsError
	}
	for index, g := range c.groups {
		if g == group {
			c.groups = append(c.groups[:index], c.groups[index+1:]...)
			return nil
		}
	}

	return NodeGroupNotFoundError
}

//--------------------------------------
// Nodes
//--------------------------------------

// Retrieves a node and its group from the cluster by id.
func (c *Cluster) GetNode(id string) (*Node, *NodeGroup) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getNode(id)
}

func (c *Cluster) getNode(id string) (*Node, *NodeGroup) {
	for _, group := range c.groups {
		if node := group.getNode(id); node != nil {
			return node, group
		}
	}
	return nil, nil
}

// Retrieves the host and port for a given node.
func (c *Cluster) GetNodeHostname(id string) (string, uint, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	node, _ := c.getNode(id)
	if node == nil {
		return "", 0, NodeNotFoundError
	}
	return node.host, node.port, nil
}

// Adds a node to an existing group in the cluster.
func (c *Cluster) AddNode(node *Node, group *NodeGroup) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Validate node.
	if node == nil {
		return NodeRequiredError
	}

	// Check if the node id exists in the cluster already.
	if n, _ := c.getNode(node.id); n != nil {
		return fmt.Errorf("Duplicate node already exists: %v", node.id)
	}

	// Find the group.
	if group == nil {
		return NodeGroupRequiredError
	}
	if group = c.getNodeGroup(group.id); group == nil {
		return NodeGroupNotFoundError
	}

	return group.addNode(node)
}

// Removes a node from a group in the cluster.
func (c *Cluster) RemoveNode(node *Node) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if node == nil {
		return NodeRequiredError
	}

	var group *NodeGroup
	if node, group = c.getNode(node.id); node == nil {
		return NodeNotFoundError
	}

	return group.removeNode(node)
}

//--------------------------------------
// Shards
//--------------------------------------

// Finds the group that the shard currently belongs to.
func (c *Cluster) GetShardOwner(index int) (*NodeGroup, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.getShardOwner(index)
}

func (c *Cluster) getShardOwner(index int) (*NodeGroup, error) {
	for _, group := range c.groups {
		if group.hasShard(index) {
			return group, nil
		}
	}
	return nil, ShardNotFoundError
}

// Moves a shard to a new group.
func (c *Cluster) MoveShard(index int, group *NodeGroup) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Find which group the shard currently belongs to.
	src, err := c.getShardOwner(index)
	if err != nil {
		return err
	}

	// Find the group the shard is moving to.
	if group == nil {
		return NodeGroupRequiredError
	}
	if group = c.getNodeGroup(group.id); group == nil {
		return NodeGroupNotFoundError
	}

	// Remove from the source group.
	if err = src.removeShard(index); err != nil {
		return err
	} else {
		// Add the shard to the destination group.
		if err = group.addShard(index); err != nil {
			// If we can't add to the new group then rollback to the old state.
			if tempErr := src.addShard(index); tempErr != nil {
				// This shouldn't ever happen but if it does then we'll shut
				// down the server because now we're missing a shard in the database.
				panic(fmt.Sprintf("skyd.Cluster: Unable to rollback shard transfer: %v", tempErr))
			}
			return err
		}
	}
	return nil
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

	return map[string]interface{}{"groups": groups}
}
