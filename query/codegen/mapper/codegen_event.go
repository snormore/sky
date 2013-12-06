package mapper

import (
	"fmt"

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

	m.builder.SetInsertPointAtEnd(m.context.AddBasicBlock(fn, "entry"))
	event := m.builder.CreateAlloca(llvm.PointerType(m.eventType, 0), "event")
	next_event := m.builder.CreateAlloca(llvm.PointerType(m.eventType, 0), "next_event")
	m.builder.CreateStore(fn.Param(0), event)
	m.builder.CreateStore(fn.Param(1), next_event)

	for _, decl := range decls {
		if decl.Id >= 0 {
			switch decl.DataType {
			case core.StringDataType:
				panic("NOT YET IMPLEMENTED: copy_permanent_variables [string]")
			case core.IntegerDataType, core.FactorDataType, core.FloatDataType, core.BooleanDataType:
				src := m.builder.CreateStructGEP(m.builder.CreateLoad(event, ""), decl.Index(), "")
				dest := m.builder.CreateStructGEP(m.builder.CreateLoad(next_event, ""), decl.Index(), "")
				m.builder.CreateStore(m.builder.CreateLoad(src, ""), dest)
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
// bool read_event(sky_event *event, void *ptr)
func (m *Mapper) codegenReadEventFunc(decls ast.VarDecls) llvm.Value {
	fntype := llvm.FunctionType(m.context.Int1Type(), []llvm.Type{llvm.PointerType(m.eventType, 0), llvm.PointerType(m.context.Int8Type(), 0)}, false)
	fn := llvm.AddFunction(m.module, "read_event", fntype)
	fn.SetFunctionCallConv(llvm.CCallConv)

	read_decls := ast.VarDecls{}
	read_labels := []llvm.BasicBlock{}
	
	entry := m.context.AddBasicBlock(fn, "entry")
	read_ts := m.context.AddBasicBlock(fn, "read_ts")
	read_map := m.context.AddBasicBlock(fn, "read_map")
	loop := m.context.AddBasicBlock(fn, "loop")
	loop_read_key := m.context.AddBasicBlock(fn, "loop_read_key")
	loop_read_value := m.context.AddBasicBlock(fn, "loop_read_value")
	loop_skip := m.context.AddBasicBlock(fn, "loop_skip")
	for _, decl := range decls {
		if decl.Id != 0 {
			read_decls = append(read_decls, decl)
			read_labels = append(read_labels, m.context.AddBasicBlock(fn, decl.Name))
		}
	}
	error := m.context.AddBasicBlock(fn, "error")
	exit := m.context.AddBasicBlock(fn, "exit")

	// entry:
	//     int64_t sz;
	//     int64_t variable_id;
	//     int64_t key_count;
	//     int64_t key_index = 0;
	m.builder.SetInsertPointAtEnd(entry)
	event := m.builder.CreateAlloca(llvm.PointerType(m.eventType, 0), "event")
	ptr := m.builder.CreateAlloca(llvm.PointerType(m.context.Int8Type(), 0), "ptr")
	sz := m.builder.CreateAlloca(m.context.Int64Type(), "sz")
	variable_id := m.builder.CreateAlloca(m.context.Int64Type(), "variable_id")
	key_count := m.builder.CreateAlloca(m.context.Int64Type(), "key_count")
	key_index := m.builder.CreateAlloca(m.context.Int64Type(), "key_index")
	m.builder.CreateStore(fn.Param(0), event)
	m.builder.CreateStore(fn.Param(1), ptr)
	m.builder.CreateStore(llvm.ConstInt(m.context.Int64Type(), 0, false), key_index)
	m.builder.CreateBr(read_ts)

	// read_ts:
	//     int64_t ts = *((int64_t*)ptr);
	//     event->timestamp = ts;
	//     ptr += 8;
	m.builder.SetInsertPointAtEnd(read_ts)
	ts_value := m.builder.CreateLoad(m.builder.CreateBitCast(m.builder.CreateLoad(ptr, ""), llvm.PointerType(m.context.Int64Type(), 0), ""), "ts_value")
	timestamp_value := m.builder.CreateLShr(ts_value, llvm.ConstInt(m.context.Int64Type(), core.SECONDS_BIT_OFFSET, false), "timestamp_value")
	event_timestamp := m.builder.CreateStructGEP(m.builder.CreateLoad(event, ""), eventTimestampElementIndex, "event_timestamp")
	m.builder.CreateStore(timestamp_value, event_timestamp)
	m.builder.CreateStore(m.builder.CreateGEP(m.builder.CreateLoad(ptr, ""), []llvm.Value{llvm.ConstInt(m.context.Int64Type(), 8, false)}, ""), ptr)
	m.builder.CreateBr(read_map)

	// read_map:
	//     key_index = 0;
	//     key_count = minipack_unpack_raw(ptr, &sz);
	//     ptr += sz;
	//     if(sz != 0) goto loop else goto error
	m.builder.SetInsertPointAtEnd(read_map)
	m.builder.CreateStore(m.builder.CreateCall(m.module.NamedFunction("minipack_unpack_map"), []llvm.Value{m.builder.CreateLoad(ptr, ""), sz}, ""), key_count)
	m.builder.CreateStore(m.builder.CreateGEP(m.builder.CreateLoad(ptr, ""), []llvm.Value{m.builder.CreateLoad(sz, "")}, ""), ptr)
	m.builder.CreateCondBr(m.builder.CreateICmp(llvm.IntNE, m.builder.CreateLoad(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop, error)

	// loop:
	//     if(key_index < key_count) goto loop_read_key else goto exit;
	m.builder.SetInsertPointAtEnd(loop)
	m.builder.CreateCondBr(m.builder.CreateICmp(llvm.IntSLT, m.builder.CreateLoad(key_index, ""), m.builder.CreateLoad(key_count, ""), ""), loop_read_key, exit)

	// loop_read_key:
	//     variable_id = minipack_unpack_int(ptr, sz)
	//     ptr += sz;
	//     if(sz != 0) goto loop_read_value else goto error;
	m.builder.SetInsertPointAtEnd(loop_read_key)
	m.builder.CreateStore(m.builder.CreateCall(m.module.NamedFunction("minipack_unpack_int"), []llvm.Value{m.builder.CreateLoad(ptr, ""), sz}, ""), variable_id)
	m.builder.CreateStore(m.builder.CreateGEP(m.builder.CreateLoad(ptr, ""), []llvm.Value{m.builder.CreateLoad(sz, "")}, ""), ptr)
	m.builder.CreateCondBr(m.builder.CreateICmp(llvm.IntNE, m.builder.CreateLoad(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop_read_value, error)

	// loop_read_value:
	//     index += 1;
	//     switch(variable_id) {
	//     ...
	//     default:
	//         goto loop_skip;
	//     }
	m.builder.SetInsertPointAtEnd(loop_read_value)
	sw := m.builder.CreateSwitch(m.builder.CreateLoad(variable_id, ""), loop_skip, len(read_decls))
	for i, decl := range read_decls {
		sw.AddCase(llvm.ConstIntFromString(m.context.Int64Type(), fmt.Sprintf("%d", decl.Id), 10), read_labels[i])
	}

	// XXX:
	//     event->XXX = minipack_unpack_XXX(ptr, sz);
	//     ptr += sz;
	//     if(sz != 0) goto loop else goto loop_skip;
	for i, decl := range read_decls {
		m.builder.SetInsertPointAtEnd(read_labels[i])

		if decl.DataType == core.StringDataType {
			panic("NOT YET IMPLEMENTED: read_event [string]")
		}

		var minipack_func_name string
		switch decl.DataType {
		case core.IntegerDataType, core.FactorDataType:
			minipack_func_name = "minipack_unpack_int"
		case core.FloatDataType:
			minipack_func_name = "minipack_unpack_double"
		case core.BooleanDataType:
			minipack_func_name = "minipack_unpack_bool"
		}

		field := m.builder.CreateStructGEP(m.builder.CreateLoad(event, ""), decl.Index(), "")
		m.builder.CreateStore(m.builder.CreateCall(m.module.NamedFunction(minipack_func_name), []llvm.Value{m.builder.CreateLoad(ptr, ""), sz}, ""), field)
		m.builder.CreateCall(m.module.NamedFunction("debug"), []llvm.Value{m.builder.CreateLoad(ptr, "")}, "")
		m.builder.CreateStore(m.builder.CreateGEP(m.builder.CreateLoad(ptr, ""), []llvm.Value{m.builder.CreateLoad(sz, "")}, ""), ptr)
		// m.createPrintfCall("sz=%d\n", m.builder.CreateLoad(sz, ""))
		m.builder.CreateCondBr(m.builder.CreateICmp(llvm.IntNE, m.builder.CreateLoad(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop, loop_skip)
	}

	// loop_skip:
	//     sz = minipack_sizeof_elem_and_data(ptr);
	//     ptr += sz;
	//     if(sz != 0) goto loop else goto error;
	m.builder.SetInsertPointAtEnd(loop_skip)
	m.builder.CreateStore(m.builder.CreateCall(m.module.NamedFunction("minipack_sizeof_elem_and_data"), []llvm.Value{m.builder.CreateLoad(ptr, "")}, ""), sz)
	m.builder.CreateStore(m.builder.CreateGEP(m.builder.CreateLoad(ptr, ""), []llvm.Value{m.builder.CreateLoad(sz, "")}, ""), ptr)
	m.builder.CreateCondBr(m.builder.CreateICmp(llvm.IntNE, m.builder.CreateLoad(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop, error)

	// error:
	//     return false;
	m.builder.SetInsertPointAtEnd(error)
	m.builder.CreateRet(llvm.ConstInt(m.context.Int1Type(), 0, false))

	// exit:
	//     return true;
	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateRet(llvm.ConstInt(m.context.Int1Type(), 1, false))

	return fn
}
