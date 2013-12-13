package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

// [codegen]
// void event_loop(sky_cursor *cursor, sky_map *result) {
// loop:
//     ...statements...
//     bool rc = cursor_next_event(cursor)
//     if(rc) goto loop else goto exit;
//
// exit:
//     return;
// }
func (m *Mapper) codegenEventLoop(node *ast.EventLoop, tbl *ast.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.hashmapType, 0)}, false)
	fn := llvm.AddFunction(m.module, "event_loop", sig)

	// Generate functions for child statements.
	var statementFns []llvm.Value
	for _, statement := range node.Statements {
		statementFn, err := m.codegenStatement(statement, tbl)
		if err != nil {
			return nilValue, err
		}
		statementFns = append(statementFns, statementFn)
	}

	entry := m.context.AddBasicBlock(fn, "entry")
	loop_condition := m.context.AddBasicBlock(fn, "loop_condition")
	loop := m.context.AddBasicBlock(fn, "loop")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor_ref := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	result_ref := m.alloca(llvm.PointerType(m.hashmapType, 0), "result")
	m.store(fn.Param(0), cursor_ref)
	m.store(fn.Param(1), result_ref)
	m.builder.CreateBr(loop_condition)

	// if(next_event->eof == false) goto loop else goto exit;
	m.builder.SetInsertPointAtEnd(loop_condition)
	event_ref := m.load_event_ref(cursor_ref)
	eof := m.load_eof(event_ref)
	m.printf("eof? %d\n", eof)
	m.condbr(m.icmp(llvm.IntEQ, eof, m.constint(0)), loop, exit)

	// ...generate...
	// if(cursor_next_event(cursor)) goto loop_condition else goto exit;
	m.builder.SetInsertPointAtEnd(loop)
	m.printf("event_loop.loop\n")
	for _, statementFn := range statementFns {
		m.call(statementFn, m.load(cursor_ref, ""), m.load(result_ref, ""))
	}
	rc := m.call("cursor_next_event", m.load(cursor_ref))
	m.printf("event_loop.next_event %d\n", rc)
	m.condbr(m.icmp(llvm.IntEQ, rc, m.constint(0)), loop_condition, exit)

	// return;
	m.builder.SetInsertPointAtEnd(exit)
	m.retvoid()

	return fn, nil
}
