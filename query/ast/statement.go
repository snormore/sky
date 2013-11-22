package ast

const (
	TypeAssignment   = "assignment"
	TypeCondition    = "condition"
	TypeDebug        = "debug"
	TypeEventLoop    = "event_loop"
	TypeExit         = "exit"
	TypeSelection    = "selection"
	TypeSessionLoop  = "session_loop"
	TypeTemporalLoop = "temporal_loop"
)

type Statement interface {
	QueryElement
	FunctionName(init bool) string
	MergeFunctionName() string
	Serialize() map[string]interface{}
	Deserialize(map[string]interface{}) error
	CodegenAggregateFunction(init bool) (string, error)
	CodegenMergeFunction(map[string]interface{}) (string, error)
	Defactorize(data interface{}) error
	Finalize(interface{}) error
	RequiresInitialization() bool
	Variables() []*Variable
	String() string
}
