package skyd

import (
	"github.com/benbjohnson/go-raft"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command removes a node group from the cluster configuration.
type RemoveNodeGroupCommand struct {
	NodeGroupId string `json:"nodeGroupId"`
}

func init() {
	raft.RegisterCommand(&RemoveNodeGroupCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewRemoveNodeGroupCommand(nodeGroupId string) *RemoveNodeGroupCommand {
	return &RemoveNodeGroupCommand{NodeGroupId: nodeGroupId}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *RemoveNodeGroupCommand) CommandName() string {
	return "group:remove"
}

func (c *RemoveNodeGroupCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)
	return server.cluster.RemoveNodeGroup(NewNodeGroup(c.NodeGroupId))
}
