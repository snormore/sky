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
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.mapType, 0)}, false)
	fn := llvm.AddFunction(m.module, "event_loop", sig)

	// Generate functions for child statements.
	var statementFns []llvm.Value
	for _, statement := range node.Statements {
		statementFn, err := m.codegen(statement, tbl)
		if err != nil {
			return nilValue, err
		}
		statementFns = append(statementFns, statementFn)
	}

	entry := m.context.AddBasicBlock(fn, "entry")
	loop := m.context.AddBasicBlock(fn, "loop")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor := m.builder.CreateAlloca(llvm.PointerType(m.cursorType, 0), "cursor")
	m.builder.CreateStore(fn.Param(0), cursor)
	result := m.builder.CreateAlloca(llvm.PointerType(m.mapType, 0), "result")
	m.builder.CreateStore(fn.Param(1), result)
	m.builder.CreateBr(loop)

	m.builder.SetInsertPointAtEnd(loop)
	for _, statementFn := range statementFns {
		m.builder.CreateCall(statementFn, []llvm.Value{m.builder.CreateLoad(cursor, ""), m.builder.CreateLoad(result, "")}, "")
	}
	rc := m.builder.CreateCall(m.module.NamedFunction("cursor_next_event"), []llvm.Value{m.builder.CreateLoad(cursor, "")}, "rc")
	m.builder.CreateCondBr(rc, loop, exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateRetVoid()

	return fn, nil
}
