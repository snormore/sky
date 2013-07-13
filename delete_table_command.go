package skyd

import (
	"github.com/benbjohnson/go-raft"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command deletes a table.
type DeleteTableCommand struct {
	Name string `json:"name"`
}

func init() {
	raft.RegisterCommand(&DeleteTableCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewDeleteTableCommand(name string) *DeleteTableCommand {
	return &DeleteTableCommand{Name: name}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *DeleteTableCommand) CommandName() string {
	return "table:delete"
}

func (c *DeleteTableCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)
	return server.DeleteTable(c.Name)
}
