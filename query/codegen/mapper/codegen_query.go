package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hashmap"
	"github.com/skydb/sky/query/codegen/minipack"
)

func (m *Mapper) codegenQuery(q *ast.Query) (llvm.Value, error) {
	tbl := ast.NewSymtable(nil)

	if _, err := m.codegenVarDecls(q.SystemVarDecls, tbl); err != nil {
		return nilValue, err
	}
	if _, err := m.codegenVarDecls(q.DeclaredVarDecls, tbl); err != nil {
		return nilValue, err
	}

	m.declare_lmdb()

	m.eventType = m.codegenEventType()
	m.cursorType = m.codegenCursorType()

	m.hashmapType = hashmap.DeclareType(m.module, m.context)
	hashmap.Declare(m.module, m.context, m.hashmapType)


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
//     if(rc) goto loop_initial_event else goto exit;
//
// loop_initial_event:
//     rc = cursor_next_event(cursor);
//     if(rc) goto loop_body else goto exit;
//
// loop_body:
//     ...generate...
//     goto loop;
//
// exit:
//     return;
// }
func (m *Mapper) codegenQueryEntryFunc(q *ast.Query, tbl *ast.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.hashmapType, 0)}, false)
	fn := llvm.AddFunction(m.module, "entry", sig)
	fn.SetFunctionCallConv(llvm.CCallConv)

	// Generate functions for child statements.
	var statementFns []llvm.Value
	for _, statement := range q.Statements {
		statementFn, err := m.codegenStatement(statement, tbl)
		if err != nil {
			return nilValue, err
		}
		statementFns = append(statementFns, statementFn)
	}

	entry := m.context.AddBasicBlock(fn, "entry")
	loop := m.context.AddBasicBlock(fn, "loop")
	loop_initial_event := m.context.AddBasicBlock(fn, "loop_initial_event")
	loop_body := m.context.AddBasicBlock(fn, "loop_body")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor_ref := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	result := m.alloca(llvm.PointerType(m.hashmapType, 0), "result")
	m.store(fn.Param(0), cursor_ref)
	m.store(fn.Param(1), result)
	m.store(m.builder.CreateMalloc(m.eventType, ""), m.structgep(m.load(cursor_ref, ""), cursorEventElementIndex, ""))
	m.store(m.builder.CreateMalloc(m.eventType, ""), m.structgep(m.load(cursor_ref, ""), cursorNextEventElementIndex, ""))
	m.printf("query (cursor=%p, event=%p, next_event=%p)\n", m.load(cursor_ref), m.load(m.structgep(m.load(cursor_ref, ""), cursorEventElementIndex, "")), m.load(m.structgep(m.load(cursor_ref, ""), cursorNextEventElementIndex, "")))
	rc := m.call("cursor_init", m.load(cursor_ref, ""))
	m.printf("\n")
	m.condbr(rc, loop_initial_event, exit)

	m.builder.SetInsertPointAtEnd(loop)
	rc = m.call("cursor_next_object", m.load(cursor_ref, ""))
	m.condbr(rc, loop_initial_event, exit)

	m.builder.SetInsertPointAtEnd(loop_initial_event)
	rc = m.call("cursor_next_event", m.load(cursor_ref))
	m.condbr(rc, loop_body, exit)

	m.builder.SetInsertPointAtEnd(loop_body)
	for _, statementFn := range statementFns {
		m.builder.CreateCall(statementFn, []llvm.Value{m.load(cursor_ref, ""), m.load(result, "")}, "")
	}
	m.printf("\n")
	m.br(loop)

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateFree(m.load(m.structgep(m.load(cursor_ref, ""), cursorEventElementIndex, ""), ""))
	m.builder.CreateFree(m.load(m.structgep(m.load(cursor_ref, ""), cursorNextEventElementIndex, ""), ""))
	m.retvoid()

	return fn, nil
}
