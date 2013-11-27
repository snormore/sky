package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	_ "github.com/skydb/sky/query/codegen"
	"github.com/skydb/sky/query/codegen/symtable"
)

// Mapper can compile a query and execute it against a cursor. The
// execution is single threaded and returns a nested map of data.
// The results can be combined using a Reducer.
type Mapper struct {
	module llvm.Module
	engine llvm.ExecutionEngine
	builder llvm.Builder
	entryFunc llvm.Value
}

// New creates a new Mapper instance.
func New(q *ast.Query) (*Mapper, error) {
	m := new(Mapper)
	m.module = llvm.NewModule("mapper")
	m.builder = llvm.NewBuilder()

	var err error
	if m.entryFunc, err = m.codegen(q, symtable.New(nil)); err != nil {
		return nil, err
	}
	if err = llvm.VerifyModule(m.module, llvm.ReturnStatusAction); err != nil {
		return nil, err
	}
	if m.engine, err = llvm.NewJITCompiler(m.module, 2); err != nil {
		return nil, err
	}

	return m, nil
}

// TODO: SetFinalizer to dispose of module, engine and builder.

// Execute runs the entry function on the execution engine.
func (m *Mapper) Execute() int {
	ret := m.engine.RunFunction(m.entryFunc, []llvm.GenericValue{})
	return int(ret.Int(false))
}

// Dump writes the LLVM IR to standard error.
func (m *Mapper) Dump() {
	m.module.Dump()
}

// Clone creates a duplicate mapper to be run independently.
func (m *Mapper) Clone() *Mapper {
	panic("NOT YET IMPLEMENTED")
}

