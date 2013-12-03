package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

const (
	eventEosElementIndex     = 0
	eventEofElementIndex = 1
	eventTimestampElementIndex = 2
)

// [codegen]
// typedef struct {
//     ...
//     struct  { void *data; int64_t sz; } string_var;
//     int64_t factor_var;
//     int64_t integer_var;
//     double  float_var;
//     bool    boolean_var;
//     ...
// } sky_event;
func (m *Mapper) codegenEventType(decls ast.VarDecls) llvm.Type {
	// Append variable declarations.
	fields := []llvm.Type{}
	for _, decl := range decls {
		var field llvm.Type

		switch decl.DataType {
		case core.StringDataType:
			field = m.context.StructType([]llvm.Type{
				m.context.Int64Type(),
				llvm.PointerType(m.context.VoidType(), 0),
			}, false)
		case core.IntegerDataType, core.FactorDataType:
			field = m.context.Int64Type()
		case core.FloatDataType:
			field = m.context.DoubleType()
		case core.BooleanDataType:
			field = m.context.Int1Type()
		}

		fields = append(fields, field)
	}

	// Return composite type.
	typ := m.context.StructCreateNamed("sky_event")
	typ.StructSetBody(fields, false)
	return typ
}

// [codegen]
// bool sky_cursor_next_event(sky_cursor *)
// bool cursor_next_event(sky_cursor *) {
//     bool rc = sky_cursor_next_event(cursor);
//     if(rc) goto cont else goto exit;
//
// cont:
//     copy_permanent_variables(cursor->event, cursor->next_event);
//     clear_transient_variables(cursor->next_event);
//     rc = read_event(cursor->next_event);
//     if(rc) goto exit else goto exit;
//
// exit:
//     return rc;
// }
func (m *Mapper) codegenCursorNextEventFunc(decls ast.VarDecls) {
	fntype := llvm.FunctionType(m.context.Int1Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false)
	fn := llvm.AddFunction(m.module, "cursor_next_event", fntype)
	cursor := fn.Param(0)
	cursor.SetName("cursor")

	copyPermanentVariablesFunc := m.codegenCopyPermanentVariablesFunc(decls)
	clearTransientVariablesFunc := m.codegenClearTransientVariablesFunc(decls)
	readEventFunc := m.codegenReadEventFunc(decls)

	phis := make([]llvm.Value, 0)
	entry := m.context.AddBasicBlock(fn, "entry")
	cont := m.context.AddBasicBlock(fn, "cont")
	exit := m.context.AddBasicBlock(fn, "exit")

	// Call C iterator.
	m.builder.SetInsertPointAtEnd(entry)
	rc := m.builder.CreateCall(m.module.NamedFunction("sky_cursor_next_event"), []llvm.Value{cursor}, "rc")
	phis = append(phis, rc)
	m.builder.CreateCondBr(rc, cont, exit)

	// Retrieve event and next_event pointers.
	m.builder.SetInsertPointAtEnd(cont)
	event := m.builder.CreateLoad(m.builder.CreateStructGEP(cursor, cursorEventElementIndex, ""), "event")
	next_event := m.builder.CreateLoad(m.builder.CreateStructGEP(cursor, cursorNextEventElementIndex, ""), "next_event")

	// Move next event to current event and read the next event.
	ptr := m.builder.CreateLoad(m.builder.CreateStructGEP(cursor, cursorPtrElementIndex, ""), "ptr")
	m.builder.CreateCall(copyPermanentVariablesFunc, []llvm.Value{event, next_event}, "")
	m.builder.CreateCall(clearTransientVariablesFunc, []llvm.Value{next_event}, "")
	rc = m.builder.CreateCall(readEventFunc, []llvm.Value{next_event, ptr}, "rc")
	phis = append(phis, rc)
	m.builder.CreateBr(exit)

	m.builder.SetInsertPointAtEnd(exit)
	phi := m.builder.CreatePHI(m.context.Int1Type(), "phi")
	phi.AddIncoming(phis, []llvm.BasicBlock{entry, cont})
	m.builder.CreateRet(phi)
}

// [codegen]
// void copy_permanent_variables(sky_event *dest, sky_event *src) {
//     ...
//     dest->field = src->field;
//     ...
// }
func (m *Mapper) codegenCopyPermanentVariablesFunc(decls ast.VarDecls) llvm.Value {
	fntype := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.eventType, 0), llvm.PointerType(m.eventType, 0)}, false)
	fn := llvm.AddFunction(m.module, "copy_permanent_variables", fntype)
	event := fn.Param(0)
	event.SetName("event")
	next_event := fn.Param(1)
	next_event.SetName("next_event")

	m.builder.SetInsertPointAtEnd(m.context.AddBasicBlock(fn, "entry"))

	for index, decl := range decls {
		if decl.Id >= 0 {
			switch decl.DataType {
			case core.StringDataType:
				panic("NOT YET IMPLEMENTED: copy_permanent_variables [string]")
			case core.IntegerDataType, core.FactorDataType, core.FloatDataType, core.BooleanDataType:
				value := m.builder.CreateLoad(m.builder.CreateStructGEP(event, index, decl.Name + "_src"), decl.Name + "_value")
				m.builder.CreateStore(value, m.builder.CreateStructGEP(event, index, decl.Name + "_dest"))
			}
		}
	}

	m.builder.CreateRetVoid()
	return fn
}

// [codegen]
// void clear_transient_variables(sky_event *event) {
//     ...
//     event->field = 0;
//     ...
// }
func (m *Mapper) codegenClearTransientVariablesFunc(decls ast.VarDecls) llvm.Value {
	fntype := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.eventType, 0)}, false)
	fn := llvm.AddFunction(m.module, "clear_transient_variables", fntype)
	event := fn.Param(0)
	event.SetName("event")

	m.builder.SetInsertPointAtEnd(m.context.AddBasicBlock(fn, "entry"))
	for index, decl := range decls {
		if decl.Id < 0 {
			switch decl.DataType {
			case core.StringDataType:
				panic("NOT YET IMPLEMENTED: clear_transient_variables [string]")
			case core.IntegerDataType, core.FactorDataType:
				m.builder.CreateStore(llvm.ConstInt(m.context.Int64Type(), 0, false), m.builder.CreateStructGEP(event, index, decl.Name))
			case core.FloatDataType:
				m.builder.CreateStore(llvm.ConstFloat(m.context.DoubleType(), 0), m.builder.CreateStructGEP(event, index, decl.Name))
			case core.BooleanDataType:
				m.builder.CreateStore(llvm.ConstInt(m.context.Int1Type(), 0, false), m.builder.CreateStructGEP(event, index, decl.Name))
			}
		}
	}

	m.builder.CreateRetVoid()
	return fn
}

// [codegen]
// bool read_event(sky_event *event, void *ptr) {
//     size_t sz;
//
// read_timestamp:
//     int64_t ts = minipack_unpack_int(ptr, &sz);
//     event->timestamp = ts;
//     ptr += sz;
//
// read_map:
//     int64_t index = 0;
//     int64_t key_count = minipack_unpack_raw(ptr, &sz);
//     if(sz == 0) return false;
//     ptr += sz;
//
// loop:
//     if(index >= key_count) goto exit;
//     index += 1;
//
//     int64_t id = minipack_unpack_int(ptr, &sz);
//     if(sz == 0) return false;
//     ptr += sz;
//     ...
//     if(id == XXX) {
//         event->field = minipack_unpack_XXX(ptr, &sz);
//         if(sz != 0) {
//             ptr += sz;
//             goto loop;
//         }
//     }
//     ...
//     ptr += minipack_sizeof_elem_and_data(ptr);
//
// exit:
//     return true;
// }
func (m *Mapper) codegenReadEventFunc(decls ast.VarDecls) llvm.Value {
	fntype := llvm.FunctionType(m.context.Int1Type(), []llvm.Type{llvm.PointerType(m.eventType, 0), llvm.PointerType(m.context.Int8Type(), 0)}, false)
	fn := llvm.AddFunction(m.module, "read_event", fntype)
	fn.SetFunctionCallConv(llvm.CCallConv)
	
	entry := m.context.AddBasicBlock(fn, "entry")
	read_ts := m.context.AddBasicBlock(fn, "read_ts")
	read_map := m.context.AddBasicBlock(fn, "read_map")
	loop := m.context.AddBasicBlock(fn, "loop")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	event := m.builder.CreateAlloca(llvm.PointerType(m.eventType, 0), "event")
	ptr := m.builder.CreateAlloca(llvm.PointerType(m.context.Int8Type(), 0), "ptr")
	sz := m.builder.CreateAlloca(m.context.Int64Type(), "sz")
	m.builder.CreateStore(fn.Param(0), event)
	m.builder.CreateStore(fn.Param(1), ptr)
	m.builder.CreateBr(read_ts)

	m.builder.SetInsertPointAtEnd(read_ts)
	ts := m.builder.CreateCall(m.module.NamedFunction("minipack_unpack_int"), []llvm.Value{m.builder.CreateLoad(ptr, ""), sz}, "ts")
	m.builder.CreateStore(ts, m.builder.CreateStructGEP(m.builder.CreateLoad(event, ""), eventTimestampElementIndex, ""))
	m.builder.CreateStore(m.builder.CreateGEP(m.builder.CreateLoad(ptr, ""), []llvm.Value{m.builder.CreateLoad(sz, "")}, ""), ptr)
	m.builder.CreateBr(read_map)

	m.builder.SetInsertPointAtEnd(read_map)
	m.builder.CreateBr(loop)

	m.builder.SetInsertPointAtEnd(loop)
	/*
	for index, decl := range decls {
		if decl.Id == 0 {
			switch decl.DataType {
			case core.StringDataType:
				panic("NOT YET IMPLEMENTED: clear_transient_variables [string]")
			case core.IntegerDataType, core.FactorDataType:
				m.builder.CreateStore(llvm.ConstInt(m.context.Int64Type(), 0, false), m.builder.CreateStructGEP(event, index, decl.Name))
			case core.FloatDataType:
				m.builder.CreateStore(llvm.ConstFloat(m.context.DoubleType(), 0), m.builder.CreateStructGEP(event, index, decl.Name))
			case core.BooleanDataType:
				m.builder.CreateStore(llvm.ConstInt(m.context.Int1Type(), 0, false), m.builder.CreateStructGEP(event, index, decl.Name))
			}
		}
	}
	*/
	m.builder.CreateBr(exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateRet(llvm.ConstInt(m.context.Int1Type(), 1, false))
	return fn
}
