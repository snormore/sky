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
	cursor := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	result := m.alloca(llvm.PointerType(m.mapType, 0), "result")
	m.store(fn.Param(0), cursor)
	m.store(fn.Param(1), result)
	m.printf("loop.0 %p\n", m.builder.CreateMalloc(m.eventType, ""))
	m.store(m.builder.CreateMalloc(m.eventType, ""), m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""))
	m.printf("loop.0 %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.store(m.builder.CreateMalloc(m.eventType, ""), m.structgep(m.load(cursor, ""), cursorNextEventElementIndex, ""))
	rc := m.builder.CreateCall(m.module.NamedFunction("cursor_init"), []llvm.Value{m.load(cursor, "")}, "rc")
	m.builder.CreateCondBr(rc, loop_body, exit)

	m.builder.SetInsertPointAtEnd(loop)
	m.printf("loop.1\n")
	rc = m.builder.CreateCall(m.module.NamedFunction("cursor_next_object"), []llvm.Value{m.load(cursor, "")}, "rc")
	m.printf("loop.1.1\n")
	m.builder.CreateCondBr(rc, loop_body, exit)

	m.builder.SetInsertPointAtEnd(loop_body)
	m.printf("loop.2\n")
	for _, statementFn := range statementFns {
		m.printf("loop.2.1\n")
		m.builder.CreateCall(statementFn, []llvm.Value{m.load(cursor, ""), m.load(result, "")}, "")
		m.printf("loop.! %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	}
	m.printf("loop.2.2\n")
	m.builder.CreateBr(loop)

	m.builder.SetInsertPointAtEnd(exit)
	m.printf("loop.3 %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.builder.CreateFree(m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.printf("loop.3.1\n")
	m.builder.CreateFree(m.load(m.structgep(m.load(cursor, ""), cursorNextEventElementIndex, ""), ""))
	m.printf("loop.3.2\n")
	m.retvoid()

	return fn, nil
}
