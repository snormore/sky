package skyd

import (
	"errors"
	"github.com/benbjohnson/go-raft"
	"time"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command deletes an event in a table.
type DeleteEventCommand struct {
	TableName string                 `json:"tableName"`
	ObjectId  string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
}

func init() {
	raft.RegisterCommand(&DeleteEventCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewDeleteEventCommand(tableName string, objectId string, timestamp time.Time) *DeleteEventCommand {
	return &DeleteEventCommand{
		TableName: tableName,
		ObjectId:  objectId,
		Timestamp: timestamp,
	}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *DeleteEventCommand) CommandName() string {
	return "event:delete"
}

func (c *DeleteEventCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)
	table, servlet, err := server.GetObjectContext(c.TableName, c.ObjectId)
	if err != nil {
		return err
	}

	if c.Timestamp.IsZero() {
		return errors.New("Timestamp required")
	}

	return servlet.DeleteEvent(table, c.ObjectId, c.Timestamp)
}
