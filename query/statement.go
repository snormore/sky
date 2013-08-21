package query

const (
	TypeAssignment = "assignment"
	TypeCondition  = "condition"
	TypeSelection  = "selection"
)

type Statement interface {
	QueryElement
	FunctionName(init bool) string
	MergeFunctionName() string
	Serialize() map[string]interface{}
	Deserialize(map[string]interface{}) error
	CodegenAggregateFunction(init bool) (string, error)
	CodegenMergeFunction() (string, error)
	Defactorize(data interface{}) error
	RequiresInitialization() bool
	String() string
}
