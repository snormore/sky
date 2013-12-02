package mapper

import (
	"github.com/axw/gollvm/llvm"
)

// [codegen]
// typedef struct {
//     sky_event *event;
//     sky_event *next_event;
// } sky_cursor;
func (m *Mapper) codegenCursorType() llvm.Type {
	typ := m.context.StructCreateNamed("sky_cursor")
	typ.StructSetBody([]llvm.Type{
		llvm.PointerType(m.eventType, 0),
		llvm.PointerType(m.eventType, 0),
	}, false)
	return typ
}

// [codegen]
// bool sky_cursor_init(sky_cursor *)
// bool sky_cursor_next_object(sky_cursor *)
func (m *Mapper) codegenCursorExternalDecl() {
	llvm.AddFunction(m.module, "sky_cursor_init", llvm.FunctionType(m.context.Int32Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false))
	llvm.AddFunction(m.module, "sky_cursor_next_object", llvm.FunctionType(m.context.Int32Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false))
	llvm.AddFunction(m.module, "sky_cursor_next_event", llvm.FunctionType(m.context.Int1Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false))
}

// [codegen]
// bool cursor_init(sky_cursor *) {
//     bool rc = sky_cursor_init(cursor);
//     return rc;
// }
func (m *Mapper) codegenCursorInitFunc() {
	sig := llvm.FunctionType(m.context.Int32Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false)
	fn := llvm.AddFunction(m.module, "cursor_init", sig)
	cursor := fn.Param(0)
	cursor.SetName("cursor")

	m.builder.SetInsertPointAtEnd(m.context.AddBasicBlock(fn, "entry"))
	rc := m.builder.CreateCall(m.module.NamedFunction("sky_cursor_init"), []llvm.Value{cursor}, "rc")
	m.builder.CreateRet(rc)
}

// [codegen]
// bool cursor_next_object(sky_cursor *) {
//     bool rc = sky_cursor_next_object(cursor);
//     return rc;
// }
func (m *Mapper) codegenCursorNextObjectFunc() {
	sig := llvm.FunctionType(m.context.Int32Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false)
	fn := llvm.AddFunction(m.module, "cursor_next_object", sig)
	cursor := fn.Param(0)
	cursor.SetName("cursor")

	m.builder.SetInsertPointAtEnd(m.context.AddBasicBlock(fn, "entry"))
	rc := m.builder.CreateCall(m.module.NamedFunction("sky_cursor_next_object"), []llvm.Value{cursor}, "rc")
	m.builder.CreateRet(rc)
}
