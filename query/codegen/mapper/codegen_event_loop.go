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
	loop := m.context.AddBasicBlock(fn, "loop")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	m.printf("event_loop.1\n")
	cursor := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	m.store(fn.Param(0), cursor)
	result := m.alloca(llvm.PointerType(m.hashmapType, 0), "result")
	m.store(fn.Param(1), result)
	m.builder.CreateBr(loop)

	m.builder.SetInsertPointAtEnd(loop)
	m.printf("event_loop.2\n")
	for _, statementFn := range statementFns {
		m.builder.CreateCall(statementFn, []llvm.Value{m.load(cursor, ""), m.load(result, "")}, "")
	}
	rc := m.builder.CreateCall(m.module.NamedFunction("cursor_next_event"), []llvm.Value{m.load(cursor, "")}, "rc")
	m.builder.CreateCondBr(rc, loop, exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.printf("event_loop.3\n")
	m.retvoid()

	return fn, nil
}
