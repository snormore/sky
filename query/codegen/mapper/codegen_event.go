package mapper

import (
	"fmt"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
)

const (
	eventEosElementIndex       = 0
	eventEofElementIndex       = 1
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
func (m *Mapper) codegenEventType() llvm.Type {
	// Append variable declarations.
	fields := []llvm.Type{}
	for _, decl := range m.decls {
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
//     MDB_val key, data;
//     goto init;
//
// init:
//     copy_permanent_variables(cursor->event, cursor->next_event);
//     clear_transient_variables(cursor->next_event);
//     goto mdb_cursor_get;
//
// mdb_cursor_get:
//     bool rc = mdb_cursor_get(cursor->lmdb_cursor, &key, &data, MDB_NEXT_DUP);
//     if(rc == 0) goto read_event else goto error;
//
// read_event:
//     rc = read_event(cursor->next_event);
//     if(rc) goto exit else goto error;
//
// error:
//     cursor->next_event->eos = true;
//     cursor->next_event->eof = true;
//     return false;
//
// exit:
//     return true;
// }
func (m *Mapper) codegenCursorNextEventFunc() {
	fntype := llvm.FunctionType(m.context.Int1Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false)
	fn := llvm.AddFunction(m.module, "cursor_next_event", fntype)

	copyPermanentVariablesFunc := m.codegenCopyPermanentVariablesFunc()
	clearTransientVariablesFunc := m.codegenClearTransientVariablesFunc()
	readEventFunc := m.codegenReadEventFunc()

	entry := m.context.AddBasicBlock(fn, "entry")
	init := m.context.AddBasicBlock(fn, "init")
	mdb_cursor_get := m.context.AddBasicBlock(fn, "mdb_cursor_get")
	read_event := m.context.AddBasicBlock(fn, "read_event")
	error_lbl := m.context.AddBasicBlock(fn, "error")
	exit := m.context.AddBasicBlock(fn, "exit")

	// Allocate stack.
	m.builder.SetInsertPointAtEnd(entry)
	cursor := m.builder.CreateAlloca(llvm.PointerType(m.cursorType, 0), "cursor")
	key := m.builder.CreateAlloca(llvm.PointerType(m.mdbValType, 0), "key")
	data := m.builder.CreateAlloca(llvm.PointerType(m.mdbValType, 0), "data")
	m.builder.CreateStore(fn.Param(0), cursor)
	event := m.builder.CreateLoad(m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorEventElementIndex, ""), "event")
	next_event := m.builder.CreateLoad(m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorNextEventElementIndex, ""), "next_event")
	m.builder.CreateBr(init)

	// Copy permanent event values and clear transient values.
	m.builder.SetInsertPointAtEnd(init)
	m.builder.CreateCall(copyPermanentVariablesFunc, []llvm.Value{event, next_event}, "")
	m.builder.CreateCall(clearTransientVariablesFunc, []llvm.Value{next_event}, "")
	m.builder.CreateBr(mdb_cursor_get)

	// Move MDB cursor forward and retrieve the new pointer.
	m.builder.SetInsertPointAtEnd(mdb_cursor_get)
	lmdb_cursor := m.builder.CreateLoad(m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorLMDBCursorElementIndex, ""), "lmdb_cursor")
	rc := m.builder.CreateCall(m.module.NamedFunction("mdb_cursor_get"), []llvm.Value{lmdb_cursor, m.builder.CreateLoad(key, ""), m.builder.CreateLoad(data, ""), llvm.ConstInt(m.context.Int64Type(), MDB_NEXT_DUP, false)}, "rc")
	m.builder.CreateCondBr(m.builder.CreateICmp(llvm.IntEQ, rc, llvm.ConstInt(m.context.Int64Type(), 0, false), ""), read_event, error_lbl)

	// Read event from pointer.
	m.builder.SetInsertPointAtEnd(read_event)
	ptr := m.builder.CreateLoad(m.builder.CreateStructGEP(m.builder.CreateLoad(cursor, ""), cursorPtrElementIndex, ""), "ptr")
	rc = m.builder.CreateCall(readEventFunc, []llvm.Value{next_event, ptr}, "rc")
	m.builder.CreateCondBr(rc, exit, error_lbl)

	m.builder.SetInsertPointAtEnd(error_lbl)
	m.builder.CreateStore(llvm.ConstInt(m.context.Int1Type(), 0, false), m.builder.CreateStructGEP(next_event, eventEosElementIndex, ""))
	m.builder.CreateStore(llvm.ConstInt(m.context.Int1Type(), 0, false), m.builder.CreateStructGEP(next_event, eventEofElementIndex, ""))
	m.builder.CreateRet(llvm.ConstInt(m.context.Int1Type(), 0, false))

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateRet(llvm.ConstInt(m.context.Int1Type(), 1, false))
}

// [codegen]
// void copy_permanent_variables(sky_event *dest, sky_event *src) {
//     ...
//     dest->field = src->field;
//     ...
// }
func (m *Mapper) codegenCopyPermanentVariablesFunc() llvm.Value {
	fntype := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.eventType, 0), llvm.PointerType(m.eventType, 0)}, false)
	fn := llvm.AddFunction(m.module, "copy_permanent_variables", fntype)

	m.builder.SetInsertPointAtEnd(m.context.AddBasicBlock(fn, "entry"))
	event := m.builder.CreateAlloca(llvm.PointerType(m.eventType, 0), "event")
	next_event := m.builder.CreateAlloca(llvm.PointerType(m.eventType, 0), "next_event")
	m.builder.CreateStore(fn.Param(0), event)
	m.builder.CreateStore(fn.Param(1), next_event)

	for _, decl := range m.decls {
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
func (m *Mapper) codegenClearTransientVariablesFunc() llvm.Value {
	fntype := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.eventType, 0)}, false)
	fn := llvm.AddFunction(m.module, "clear_transient_variables", fntype)
	event := fn.Param(0)
	event.SetName("event")

	m.builder.SetInsertPointAtEnd(m.context.AddBasicBlock(fn, "entry"))
	for index, decl := range m.decls {
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
func (m *Mapper) codegenReadEventFunc() llvm.Value {
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
	for _, decl := range m.decls {
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
