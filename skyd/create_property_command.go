package skyd

import (
	"github.com/benbjohnson/go-raft"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command creates a property on a table.
type CreatePropertyCommand struct {
	TableName string `json:"tableName"`
	Name      string `json:"name"`
	Transient bool   `json:"transient"`
	DataType  string `json:"dataType"`
}

func init() {
	raft.RegisterCommand(&CreatePropertyCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewCreatePropertyCommand(tableName string, name string, transient bool, dataType string) *CreatePropertyCommand {
	return &CreatePropertyCommand{
		TableName: tableName,
		Name:      name,
		Transient: transient,
		DataType:  dataType,
	}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *CreatePropertyCommand) CommandName() string {
	return "property:create"
}

func (c *CreatePropertyCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)

	table, err := server.OpenTable(c.TableName)
	if err != nil {
		return err
	}

	_, err = table.CreateProperty(c.Name, c.Transient, c.DataType)
	return err
}
