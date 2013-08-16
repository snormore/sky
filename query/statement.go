package query

const (
	TypeCondition = "condition"
	TypeSelection = "selection"
)

type Statement interface {
	FunctionName(init bool) string
	MergeFunctionName() string
	GetStatements() Statements
	Serialize() map[string]interface{}
	Deserialize(map[string]interface{}) error
	CodegenAggregateFunction(init bool) (string, error)
	CodegenMergeFunction() (string, error)
	Defactorize(data interface{}) error
	RequiresInitialization() bool
	String() string
}

