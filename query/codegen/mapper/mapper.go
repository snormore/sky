package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	_ "github.com/skydb/sky/query/codegen"
)

// Mapper can compile a query and execute it against a cursor. The
// execution is single threaded and returns a nested map of data.
// The results can be combined using a Reducer.
type Mapper struct {
	module llvm.Module
	engine llvm.ExecutionEngine
}

// New creates a new Mapper instance.
func New(q *ast.Query) (*Mapper, error) {
	m := new(Mapper)
	m.codegen(q)
	return m, nil
}

// Execute runs the entry function on the execution engine.
func (m *Mapper) Execute() int {
	fn := llvm.AddFunction(m.module, "entry", llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false))
	fn.SetFunctionCallConv(llvm.CCallConv)
	return 0
}

// Clone creates a duplicate mapper to be run independently.
func (m *Mapper) Clone() *Mapper {
	panic("NOT YET IMPLEMENTED")
}

// codegen generates the execution engine and module for the mapper.
func (m *Mapper) codegen(q *ast.Query) error {
	m.module = llvm.NewModule("mapper")
	entryFunc := llvm.AddFunction(m.module, "entry", llvm.FunctionType(llvm.Int32Type(), []llvm.Type{}, false))
	entryFunc.SetFunctionCallConv(llvm.CCallConv)
	return nil
}

