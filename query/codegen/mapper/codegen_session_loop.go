package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

// [codegen]
// void session_loop(sky_cursor *cursor, sky_map *result) {
// init:
//     int64_t prev_session_idle_time = cursor->session_idle_time;
//     cursor->session_idle_time = <SESSION_IDLE>;
//
// loop:
//     ...statements...
//     bool rc = sky_cursor_next_session(cursor)
//     if(rc) goto loop else goto exit;
//
// exit:
//     cursor->session_idle_time = prev_session_idle_time;
//     return;
// }
func (m *Mapper) codegenSessionLoop(node *ast.SessionLoop, tbl *ast.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.hashmapType, 0)}, false)
	fn := llvm.AddFunction(m.module, "session_loop", sig)

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
	init := m.context.AddBasicBlock(fn, "init")
	loop_condition := m.context.AddBasicBlock(fn, "loop_condition")
	loop := m.context.AddBasicBlock(fn, "loop")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor_ref := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	result_ref := m.alloca(llvm.PointerType(m.hashmapType, 0), "result")
	m.store(fn.Param(0), cursor_ref)
	m.store(fn.Param(1), result_ref)
	m.br(init)

	// int64_t prev_session_idle_time = cursor->session_idle_time;
	// cursor->session_idle_time = <SESSION_IDLE>;
	m.builder.SetInsertPointAtEnd(init)
	prevSessionIdleTime := m.load(m.structgep(m.load(cursor_ref), cursorSessionIdleTimeElementIndex))
	m.store(m.constint(node.IdleDuration), m.structgep(m.load(cursor_ref), cursorSessionIdleTimeElementIndex))
	m.printf("session_loop.init %d\n", m.load(m.structgep(m.load(cursor_ref), cursorSessionIdleTimeElementIndex)))
	m.br(loop_condition)

	// if(next_event->eof == false) goto loop else goto exit;
	m.builder.SetInsertPointAtEnd(loop_condition)
	event_ref := m.event_ref(cursor_ref)
	eof := m.load_eof(event_ref)
	m.condbr(m.icmp(llvm.IntEQ, eof, m.constint(0)), loop, exit)

	// ...generate...
	// cursor->session_wait = 0;
	// rc = cursor_next_event(cursor);
	// if (rc == 0) goto loop_condition else goto exit;
	m.builder.SetInsertPointAtEnd(loop)
	m.printf("session_loop.loop\n")
	for _, statementFn := range statementFns {
		m.call(statementFn, m.load(cursor_ref, ""), m.load(result_ref, ""))
	}
	m.store(m.constint(0), m.structgep(m.load(cursor_ref), cursorSessionWaitElementIndex))
	rc := m.call("cursor_next_event", m.load(cursor_ref))
	m.printf("session_loop.next_event %d\n", rc)
	m.condbr(m.icmp(llvm.IntEQ, rc, m.constint(0)), loop_condition, exit)

	// cursor->session_idle_time = prev_session_idle_time;
	// return;
	m.builder.SetInsertPointAtEnd(exit)
	m.store(prevSessionIdleTime, m.structgep(m.load(cursor_ref), cursorSessionIdleTimeElementIndex))
	m.retvoid()

	return fn, nil
}
