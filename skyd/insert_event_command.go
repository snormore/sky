package skyd

import (
	"github.com/benbjohnson/go-raft"
	"time"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command inserts an event into a table.
type InsertEventCommand struct {
	TableName string `json:"tableName"`
	ObjectId      string `json:"id"`
	Timestamp time.Time   `json:"timestamp"`
	Data  map[string]interface{} `json:"data"`
	Method string `json:"method"`
}

func init() {
	raft.RegisterCommand(&InsertEventCommand{})
}

//------------------------------------------------------------------------------
//
// Methods
//
//------------------------------------------------------------------------------

//--------------------------------------
// Constructor
//--------------------------------------

func NewInsertEventCommand(tableName string, objectId string, timestamp time.Time, data map[string]interface{}, method string) *InsertEventCommand {
	return &InsertEventCommand{
		TableName: tableName,
		ObjectId:      objectId,
		Timestamp:      timestamp,
		Data:      data,
		Method:      method,
	}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *InsertEventCommand) CommandName() string {
	return "event:insert"
}

func (c *InsertEventCommand) Apply(raftServer *raft.Server) error {
	//server := raftServer.Context().(*Server)
	return nil
}

