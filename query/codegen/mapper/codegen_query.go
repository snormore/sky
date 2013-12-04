package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/minipack"
	"github.com/skydb/sky/query/codegen/symtable"
)

func (m *Mapper) codegenQuery(q *ast.Query, tbl *symtable.Symtable) (llvm.Value, error) {
	// Generate "event" struct type.
	decls, err := q.VarDecls()
	if err != nil {
		return nilValue, err
	}
	m.eventType = m.codegenEventType(decls)
	m.cursorType = m.codegenCursorType()
	m.mapType = m.context.StructCreateNamed("sky_map")

	llvm.AddFunction(m.module, "debug", llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.context.Int8Type(), 0)}, false))
	m.codegenCursorExternalDecl()

	minipack.Declare_unpack_int(m.module, m.context)
	minipack.Declare_unpack_double(m.module, m.context)
	minipack.Declare_unpack_bool(m.module, m.context)
	minipack.Declare_unpack_raw(m.module, m.context)
	minipack.Declare_unpack_map(m.module, m.context)
	minipack.Declare_sizeof_elem_and_data(m.module, m.context)

	m.codegenCursorInitFunc()
	m.codegenCursorNextObjectFunc()
	m.codegenCursorNextEventFunc(decls)

	// Generate the entry function.
	return m.codegenQueryEntryFunc(q, tbl)
}

// [codegen]
// int32_t entry(sky_cursor *cursor, sky_map *result) {
//     int32_t rc = cursor_init(cursor);
//     if(rc != 0) goto exit;
//
// loop:
//     rc = cursor_next_object(cursor);
//     if(rc != 0) goto exit;
//     goto loop
//
// exit:
//     return rc;
// }
func (m *Mapper) codegenQueryEntryFunc(q *ast.Query, tbl *symtable.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.Int32Type(), []llvm.Type{
		llvm.PointerType(m.cursorType, 0),
		llvm.PointerType(m.mapType, 0),
	}, false)
	fn := llvm.AddFunction(m.module, "entry", sig)
	fn.SetFunctionCallConv(llvm.CCallConv)
	cursor := fn.Param(0)
	cursor.SetName("cursor")
	result := fn.Param(1)
	result.SetName("result")

	entry := m.context.AddBasicBlock(fn, "entry")
	//loop := m.context.AddBasicBlock(fn, "loop")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	m.builder.CreateCall(m.module.NamedFunction("cursor_init"), []llvm.Value{cursor}, "rc")
	m.builder.CreateBr(exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateRet(llvm.ConstInt(m.context.Int32Type(), 12, false))
	return fn, nil
}
