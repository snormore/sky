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

// This command deletes a property on a table.
type DeletePropertyCommand struct {
	TableName string `json:"tableName"`
	Name      string `json:"name"`
}

func init() {
	raft.RegisterCommand(&DeletePropertyCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewDeletePropertyCommand(tableName string, name string) *DeletePropertyCommand {
	return &DeletePropertyCommand{
		TableName: tableName,
		Name:      name,
	}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *DeletePropertyCommand) CommandName() string {
	return "property:delete"
}

func (c *DeletePropertyCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)
	table, err := server.OpenTable(c.TableName)
	if err != nil {
		return err
	}
	// Retrieve property.
	property, err := table.GetPropertyByName(c.Name)
	if err != nil {
		return err
	} else if property == nil {
		return fmt.Errorf("Property not found: %s", c.Name)
	}

	// Delete property and save property file.
	table.DeleteProperty(property)
	return table.SavePropertyFile()
}
