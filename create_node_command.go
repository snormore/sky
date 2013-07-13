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

// This command creates a node in the cluster configuration.
type CreateNodeCommand struct {
	NodeId      string `json:"nodeId"`
	NodeGroupId string `json:"nodeGroupId"`
	Host        string `json:"host"`
	Port        uint   `json:"port"`
}

func init() {
	raft.RegisterCommand(&CreateNodeCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewCreateNodeCommand(nodeId string, nodeGroupId string, host string, port uint) *CreateNodeCommand {
	return &CreateNodeCommand{
		NodeId:      nodeId,
		NodeGroupId: nodeGroupId,
		Host:        host,
		Port:        port,
	}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *CreateNodeCommand) CommandName() string {
	return "node:create"
}

func (c *CreateNodeCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)

	// TODO: Validate host & port.

	// Create node and add it to the cluster.
	node := NewNode(c.NodeId, c.Host, c.Port)
	group := NewNodeGroup(c.NodeGroupId)
	err := server.cluster.AddNode(node, group)
	if err != nil {
		return err
	}

	// Add the node as a peer.
	if err := server.clusterRaftServer.AddPeer(node.id); err != nil {
		// This should never error so if it does then our cluster configuration
		// and list of peers are out of sync and we should just crash.
		panic(fmt.Sprintf("skyd: %v", err))
	}

	return nil
}
