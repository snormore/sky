package skyd

import (
	"errors"
	"github.com/benbjohnson/go-raft"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command creates a table.
type CreateTableCommand struct {
	Name      string `json:"name"`
}

func init() {
	raft.RegisterCommand(&CreateTableCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewCreateTableCommand(name string) *CreateTableCommand {
	return &CreateTableCommand{Name:name}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *CreateTableCommand) CommandName() string {
	return "table:create"
}

func (c *CreateTableCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)

	// Retrieve table parameters.
	if c.Name == "" {
		return errors.New("Table name required.")
	}

	// Return an error if the table already exists.
	table, err := server.OpenTable(c.Name)
	if table != nil {
		return errors.New("Table already exists.")
	}

	// Otherwise create it.
	table = NewTable(c.Name, server.TablePath(c.Name))
	if err = table.Create(); err != nil {
		return err
	}

	return nil
}
