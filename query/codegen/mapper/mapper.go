package mapper

// #include <stdlib.h>
import "C"

import (
	"runtime"
	"unsafe"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	_ "github.com/skydb/sky/query/codegen"
	"github.com/skydb/sky/query/codegen/hashmap"
	"github.com/szferi/gomdb"
)

// Mapper can compile a query and execute it against a cursor. The
// execution is single threaded and returns a nested map of data.
// The results can be combined using a Reducer.
type Mapper struct {
	context llvm.Context
	module  llvm.Module
	engine  llvm.ExecutionEngine
	builder llvm.Builder

	decls    ast.VarDecls

	cursorType    llvm.Type
	eventType     llvm.Type
	hashmapType   llvm.Type
	mdbCursorType llvm.Type
	mdbValType    llvm.Type

	entryFunc llvm.Value
}

// New creates a new Mapper instance.
func New(q *ast.Query) (*Mapper, error) {
	m := new(Mapper)
	runtime.SetFinalizer(m, finalize)

	m.context = llvm.NewContext()
	m.module = m.context.NewModule("mapper")
	m.builder = llvm.NewBuilder()

	var err error
	if m.decls, err = q.VarDecls(); err != nil {
		return nil, err
	}
	if m.entryFunc, err = m.codegenQuery(q); err != nil {
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

// finalize cleans up resources after the mapper goes out of scope.
func finalize(m *Mapper) {
	m.builder.Dispose()
	m.engine.Dispose()
}

// Execute runs the entry function on the execution engine.
func (m *Mapper) Execute(lmdb_cursor *mdb.Cursor, prefix string, result *hashmap.Hashmap) error {
	cursor := sky_cursor_new(lmdb_cursor, prefix)
	defer sky_cursor_free(cursor)

	m.engine.RunFunction(m.entryFunc, []llvm.GenericValue{
		llvm.NewGenericValueFromPointer(unsafe.Pointer(cursor)),
		llvm.NewGenericValueFromPointer(unsafe.Pointer(result)),
	})
	return nil
}

// Execute a given function and returns the result.
func (m *Mapper) ExecuteFunction(name string, values []llvm.GenericValue) llvm.GenericValue {
	fn := m.module.NamedFunction(name)
	return m.engine.RunFunction(fn, values)
}

// Dump writes the LLVM IR to STDERR.
func (m *Mapper) Dump() {
	m.module.Dump()
}

// Clone creates a duplicate mapper to be run independently.
func (m *Mapper) Clone() *Mapper {
	panic("NOT YET IMPLEMENTED")
}
