package skyd

import (
	"github.com/benbjohnson/go-raft"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command creates a marker in the log that indicates where initialization
// is complete. A server can only join a cluster if it's last log entry is an
// init command. This command is a noop.
type InitCommand struct {
}

func init() {
	raft.RegisterCommand(&InitCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewInitCommand() *InitCommand {
	return &InitCommand{}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *InitCommand) CommandName() string {
	return "init"
}

func (c *InitCommand) Apply(raftServer *raft.Server) error {
	return nil
}
