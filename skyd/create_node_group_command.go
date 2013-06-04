package skyd

import (
	"github.com/benbjohnson/go-raft"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command creates a node group in the cluster configuration.
type CreateNodeGroupCommand struct {
	NodeGroupId string `json:"nodeGroupId"`
}

func init() {
	raft.RegisterCommand(&CreateNodeGroupCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewCreateNodeGroupCommand(nodeGroupId string) *CreateNodeGroupCommand {
	return &CreateNodeGroupCommand{NodeGroupId: nodeGroupId}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *CreateNodeGroupCommand) CommandName() string {
	return "group:create"
}

func (c *CreateNodeGroupCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)
	group := NewNodeGroup(c.NodeGroupId)
	return server.cluster.AddNodeGroup(group)
}
