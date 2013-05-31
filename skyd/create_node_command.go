package skyd

import (
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

//------------------------------------------------------------------------------
//
// Constructor
//
//------------------------------------------------------------------------------

func NewCreateNodeCommand(nodeId string, nodeGroupId string, host string, port uint) *CreateNodeCommand {
	return &CreateNodeCommand{
		NodeId:      nodeId,
		NodeGroupId: nodeGroupId,
		Host:        host,
		Port:        port,
	}
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

func (c *CreateNodeCommand) CommandName() string {
	return "node:create"
}

func (c *CreateNodeCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)

	warn("cnc.apply.1: %v", c)
	// Locate the group to add to.
	group := server.cluster.GetNodeGroup(c.NodeGroupId)
	if group == nil {
		return NodeGroupNotFoundError
	}

	warn("cnc.apply.2")
	// Create node and add it to the cluster.
	node := NewNode(c.NodeId, c.Host, c.Port)
	warn("cnc.apply.3")
	return server.cluster.AddNode(node, group)
}
