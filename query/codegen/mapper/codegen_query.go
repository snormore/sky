package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/minipack"
)

func (m *Mapper) codegenQuery(q *ast.Query, tbl *ast.Symtable) (llvm.Value, error) {
	m.declare_lmdb()

	m.eventType = m.codegenEventType()
	m.cursorType = m.codegenCursorType()
	m.mapType = m.context.StructCreateNamed("sky_map")

	llvm.AddFunction(m.module, "debug", llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.context.Int8Type(), 0)}, false))
	m.codegenCursorExternalDecl()

	llvm.AddFunction(m.module, "printf", llvm.FunctionType(m.context.Int32Type(), []llvm.Type{}, true))

	minipack.Declare_unpack_int(m.module, m.context)
	minipack.Declare_unpack_double(m.module, m.context)
	minipack.Declare_unpack_bool(m.module, m.context)
	minipack.Declare_unpack_raw(m.module, m.context)
	minipack.Declare_unpack_map(m.module, m.context)
	minipack.Declare_sizeof_elem_and_data(m.module, m.context)

	m.codegenCursorInitFunc()
	m.codegenCursorNextObjectFunc()
	m.codegenCursorNextEventFunc()

	// Generate the entry function.
	return m.codegenQueryEntryFunc(q, tbl)
}

// [codegen]
// int32_t entry(sky_cursor *cursor, sky_map *result) {
//     cursor->event = malloc();
//     cursor->next_event = malloc();
//     int32_t rc = cursor_init(cursor);
//     if(rc) goto loop_body else goto exit;
//
// loop:
//     rc = cursor_next_object(cursor);
//     if(rc) goto loop_body else goto exit;
//
// loop_body:
//     query(cursor, result);
//     goto loop;
//
// exit:
//     return;
// }
func (m *Mapper) codegenQueryEntryFunc(q *ast.Query, tbl *ast.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.mapType, 0)}, false)
	fn := llvm.AddFunction(m.module, "entry", sig)
	fn.SetFunctionCallConv(llvm.CCallConv)

	// Generate functions for child statements.
	var statementFns []llvm.Value
	for _, statement := range q.Statements {
		statementFn, err := m.codegen(statement, tbl)
		if err != nil {
			return nilValue, err
		}
		statementFns = append(statementFns, statementFn)
	}

	entry := m.context.AddBasicBlock(fn, "entry")
	loop := m.context.AddBasicBlock(fn, "loop")
	loop_body := m.context.AddBasicBlock(fn, "loop_body")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor := m.builder.CreateAlloca(llvm.PointerType(m.cursorType, 0), "cursor")
	result := m.builder.CreateAlloca(llvm.PointerType(m.mapType, 0), "result")
	m.builder.CreateStore(fn.Param(0), cursor)
	m.builder.CreateStore(fn.Param(1), result)
	m.builder.CreateStore(m.builder.CreateMalloc(m.eventType, ""), m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorEventElementIndex, ""))
	m.builder.CreateStore(m.builder.CreateMalloc(m.eventType, ""), m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorNextEventElementIndex, ""))
	rc := m.builder.CreateCall(m.module.NamedFunction("cursor_init"), []llvm.Value{m.builder.CreateLoad(cursor, "")}, "rc")
	m.builder.CreateCondBr(rc, loop_body, exit)

	m.builder.SetInsertPointAtEnd(loop)
	rc = m.builder.CreateCall(m.module.NamedFunction("cursor_next_object"), []llvm.Value{m.builder.CreateLoad(cursor, "")}, "rc")
	m.builder.CreateCondBr(rc, loop_body, exit)

	m.builder.SetInsertPointAtEnd(loop_body)
	for _, statementFn := range statementFns {
		m.builder.CreateCall(statementFn, []llvm.Value{m.builder.CreateLoad(cursor, ""), m.builder.CreateLoad(result, "")}, "")
	}
	m.builder.CreateBr(loop)

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateFree(m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorEventElementIndex, ""))
	m.builder.CreateFree(m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorNextEventElementIndex, ""))
	m.builder.CreateRetVoid()

	return fn, nil
}
