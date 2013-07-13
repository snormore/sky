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

// This command updates a property on a table.
type UpdatePropertyCommand struct {
	TableName    string `json:"tableName"`
	OriginalName string `json:"originalName"`
	Name         string `json:"name"`
}

func init() {
	raft.RegisterCommand(&UpdatePropertyCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewUpdatePropertyCommand(tableName string, originalName string, name string) *UpdatePropertyCommand {
	return &UpdatePropertyCommand{
		TableName:    tableName,
		OriginalName: originalName,
		Name:         name,
	}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *UpdatePropertyCommand) CommandName() string {
	return "property:update"
}

func (c *UpdatePropertyCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)
	table, err := server.OpenTable(c.TableName)
	if err != nil {
		return err
	}

	// Retrieve property.
	property, err := table.GetPropertyByName(c.OriginalName)
	if err != nil {
		return err
	} else if property == nil {
		return fmt.Errorf("Property does not exist: %v", c.OriginalName)
	}

	// Update property and save property file.
	property.Name = c.Name
	return table.SavePropertyFile()
}
