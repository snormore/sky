package skyd

import (
	"fmt"
	"github.com/benbjohnson/go-raft"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command removes a node from the cluster configuration.
type RemoveNodeCommand struct {
	NodeId string `json:"nodeId"`
}

func init() {
	raft.RegisterCommand(&RemoveNodeCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewRemoveNodeCommand(nodeId string) *RemoveNodeCommand {
	return &RemoveNodeCommand{NodeId: nodeId}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *RemoveNodeCommand) CommandName() string {
	return "node:remove"
}

func (c *RemoveNodeCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)

	// Don't allow the last node to be removed from a group.
	node, group := server.cluster.GetNode(c.NodeId)
	if node == nil {
		return fmt.Errorf("Node not found: %s", c.NodeId)
	} else if len(group.nodes) == 1 {
		return fmt.Errorf("Cannot remove last node from group: %s/%s", node.id, group.id)
	}

	// If the node being removed is this server then clear out the data
	// and shutdown.
	if server.name == c.NodeId {
		// TODO: Clear all data.
		go func() { server.Shutdown() }()
	} else {
		// Remove the node from the set of peers.
		if err := server.clusterRaftServer.RemovePeer(c.NodeId); err != nil {
			return err
		}

		// Remove node from the cluster.
		node := NewNode(c.NodeId, "", 0)
		if err := server.cluster.RemoveNode(node); err != nil {
			return err
		}
	}
	return nil
}
