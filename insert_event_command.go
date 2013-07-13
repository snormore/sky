package skyd

import (
	"errors"
	"github.com/benbjohnson/go-raft"
	"time"
)

//------------------------------------------------------------------------------
//
// Globals
//
//------------------------------------------------------------------------------

const (
	ReplaceMethod = "replace"
	MergeMethod = "merge"
)

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// This command inserts an event into a table.
type InsertEventCommand struct {
	TableName string                 `json:"tableName"`
	ObjectId  string                 `json:"id"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Method    string                 `json:"method"`
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
		ObjectId:  objectId,
		Timestamp: timestamp,
		Data:      data,
		Method:    method,
	}
}

//--------------------------------------
// Command
//--------------------------------------

func (c *InsertEventCommand) CommandName() string {
	return "event:insert"
}

func (c *InsertEventCommand) Apply(raftServer *raft.Server) error {
	server := raftServer.Context().(*Server)
	table, servlet, err := server.GetObjectContext(c.TableName, c.ObjectId)
	if err != nil {
		return err
	}

	if c.Timestamp.IsZero() {
		return errors.New("Timestamp required")
	}

	// Build event, factorize it and insert it.
	event := &Event{Timestamp:c.Timestamp}
	event.Data, err = table.NormalizeMap(c.Data)
	if err = table.FactorizeEvent(event, server.factors, true); err != nil {
		return err
	}

	replace := (c.Method == ReplaceMethod)
	err = servlet.PutEvent(table, c.ObjectId, event, replace)
	return err
}
