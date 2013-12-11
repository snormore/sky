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
	m.printf("next_event.1\n")
	cursor := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	key := m.alloca(m.mdbValType, "key")
	data := m.alloca(m.mdbValType, "data")
	m.store(fn.Param(0), cursor)
	event := m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), "event")
	next_event := m.load(m.structgep(m.load(cursor, ""), cursorNextEventElementIndex, ""), "next_event")
	m.printf("next_event.2\n")
	m.printf("next_event >>> %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.br(init)

	// Copy permanent event values and clear transient values.
	m.builder.SetInsertPointAtEnd(init)
	m.builder.CreateCall(copyPermanentVariablesFunc, []llvm.Value{event, next_event}, "")
	m.builder.CreateCall(clearTransientVariablesFunc, []llvm.Value{next_event}, "")
	m.printf("next_event.3\n")
	m.printf("next_event >>> %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.br(mdb_cursor_get)

	// Move MDB cursor forward and retrieve the new pointer.
	m.builder.SetInsertPointAtEnd(mdb_cursor_get)
	lmdb_cursor := m.load(m.structgep(m.load(cursor, ""), cursorLMDBCursorElementIndex, ""), "lmdb_cursor")
	m.printf("next_event >>> %p | %p | %p\n",
		m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""),
		m.load(cursor, ""),
		lmdb_cursor)
	rc := m.builder.CreateCall(m.module.NamedFunction("mdb_cursor_get"), []llvm.Value{lmdb_cursor, key, data, llvm.ConstInt(m.context.Int64Type(), MDB_NEXT_DUP, false)}, "rc")
	m.printf("next_event.4\n")
	m.printf("next_event >>> %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.condbr(m.builder.CreateICmp(llvm.IntEQ, rc, llvm.ConstInt(m.context.Int64Type(), 0, false), ""), read_event, error_lbl)

	// Read event from pointer.
	m.builder.SetInsertPointAtEnd(read_event)
	m.printf("next_event.4.1\n")
	ptr := m.load(m.structgep(data, 1, ""), "ptr")
	m.printf("next_event.4.2 %p\n", ptr)
	rc = m.builder.CreateCall(readEventFunc, []llvm.Value{next_event, ptr}, "rc")
	m.printf("next_event.5 %d\n", rc)
	m.printf("next_event >>> %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.condbr(rc, exit, error_lbl)

	m.builder.SetInsertPointAtEnd(error_lbl)
	m.printf("next_event.6.1\n")
	m.store(llvm.ConstInt(m.context.Int1Type(), 0, false), m.structgep(next_event, eventEosElementIndex, ""))
	m.printf("next_event.6.2\n")
	m.store(llvm.ConstInt(m.context.Int1Type(), 0, false), m.structgep(next_event, eventEofElementIndex, ""))
	m.printf("next_event.6.3\n")
	m.printf("next_event >>> %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.ret(llvm.ConstInt(m.context.Int1Type(), 0, false))

	m.builder.SetInsertPointAtEnd(exit)
	m.printf("next_event.7\n")
	m.printf("next_event >>> %p\n", m.load(m.structgep(m.load(cursor, ""), cursorEventElementIndex, ""), ""))
	m.ret(llvm.ConstInt(m.context.Int1Type(), 1, false))
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
	event := m.alloca(llvm.PointerType(m.eventType, 0), "event")
	next_event := m.alloca(llvm.PointerType(m.eventType, 0), "next_event")
	m.store(fn.Param(0), event)
	m.store(fn.Param(1), next_event)

	for _, decl := range m.decls {
		if decl.Id >= 0 {
			switch decl.DataType {
			case core.StringDataType:
				panic("NOT YET IMPLEMENTED: copy_permanent_variables [string]")
			case core.IntegerDataType, core.FactorDataType, core.FloatDataType, core.BooleanDataType:
				src := m.structgep(m.load(event, ""), decl.Index(), "")
				dest := m.structgep(m.load(next_event, ""), decl.Index(), "")
				m.store(m.load(src, ""), dest)
			}
		}
	}

	m.retvoid()
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
				m.store(llvm.ConstInt(m.context.Int64Type(), 0, false), m.structgep(event, index, decl.Name))
			case core.FloatDataType:
				m.store(llvm.ConstFloat(m.context.DoubleType(), 0), m.structgep(event, index, decl.Name))
			case core.BooleanDataType:
				m.store(llvm.ConstInt(m.context.Int1Type(), 0, false), m.structgep(event, index, decl.Name))
			}
		}
	}

	m.retvoid()
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
	m.printf("read.1\n")
	event := m.alloca(llvm.PointerType(m.eventType, 0), "event")
	ptr := m.alloca(llvm.PointerType(m.context.Int8Type(), 0), "ptr")
	sz := m.alloca(m.context.Int64Type(), "sz")
	variable_id := m.alloca(m.context.Int64Type(), "variable_id")
	key_count := m.alloca(m.context.Int64Type(), "key_count")
	key_index := m.alloca(m.context.Int64Type(), "key_index")
	m.store(fn.Param(0), event)
	m.store(fn.Param(1), ptr)
	m.store(llvm.ConstInt(m.context.Int64Type(), 0, false), key_index)
	m.printf("read.2\n")
	m.br(read_ts)

	// read_ts:
	//     int64_t ts = *((int64_t*)ptr);
	//     event->timestamp = ts;
	//     ptr += 8;
	m.builder.SetInsertPointAtEnd(read_ts)
	m.printf("read.3.1 %p\n", m.load(ptr, ""))
	ts_value := m.load(m.builder.CreateBitCast(m.load(ptr, ""), llvm.PointerType(m.context.Int64Type(), 0), ""), "ts_value")
	m.printf("read.3.2\n")
	timestamp_value := m.builder.CreateLShr(ts_value, llvm.ConstInt(m.context.Int64Type(), core.SECONDS_BIT_OFFSET, false), "timestamp_value")
	m.printf("read.3.3\n")
	event_timestamp := m.structgep(m.load(event, ""), eventTimestampElementIndex, "event_timestamp")
	m.printf("read.3.4\n")
	m.store(timestamp_value, event_timestamp)
	m.printf("read.3.5\n")
	m.store(m.builder.CreateGEP(m.load(ptr, ""), []llvm.Value{llvm.ConstInt(m.context.Int64Type(), 8, false)}, ""), ptr)
	m.printf("read.3.6\n")
	m.br(read_map)

	// read_map:
	//     key_index = 0;
	//     key_count = minipack_unpack_raw(ptr, &sz);
	//     ptr += sz;
	//     if(sz != 0) goto loop else goto error
	m.builder.SetInsertPointAtEnd(read_map)
	m.store(m.builder.CreateCall(m.module.NamedFunction("minipack_unpack_map"), []llvm.Value{m.load(ptr, ""), sz}, ""), key_count)
	m.store(m.builder.CreateGEP(m.load(ptr, ""), []llvm.Value{m.load(sz, "")}, ""), ptr)
	m.printf("read.4 (key_count=%d)\n", m.load(key_count))
	m.condbr(m.builder.CreateICmp(llvm.IntNE, m.load(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop, error)

	// loop:
	//     if(key_index < key_count) goto loop_read_key else goto exit;
	m.builder.SetInsertPointAtEnd(loop)
	m.printf("read.5\n")
	m.condbr(m.builder.CreateICmp(llvm.IntSLT, m.load(key_index, ""), m.load(key_count, ""), ""), loop_read_key, exit)

	// loop_read_key:
	//     variable_id = minipack_unpack_int(ptr, sz)
	//     ptr += sz;
	//     if(sz != 0) goto loop_read_value else goto error;
	m.builder.SetInsertPointAtEnd(loop_read_key)
	m.store(m.builder.CreateCall(m.module.NamedFunction("minipack_unpack_int"), []llvm.Value{m.load(ptr, ""), sz}, ""), variable_id)
	m.store(m.builder.CreateAdd(m.load(key_index), llvm.ConstInt(m.context.Int64Type(), 1, false), ""), key_index)
	m.store(m.builder.CreateGEP(m.load(ptr, ""), []llvm.Value{m.load(sz, "")}, ""), ptr)
	m.printf("read.6 %d | (next=%p)\n", m.load(variable_id), m.load(ptr))
	m.condbr(m.builder.CreateICmp(llvm.IntNE, m.load(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop_read_value, error)

	// loop_read_value:
	//     index += 1;
	//     switch(variable_id) {
	//     ...
	//     default:
	//         goto loop_skip;
	//     }
	m.builder.SetInsertPointAtEnd(loop_read_value)
	sw := m.builder.CreateSwitch(m.load(variable_id, ""), loop_skip, len(read_decls))
	for i, decl := range read_decls {
		sw.AddCase(llvm.ConstIntFromString(m.context.Int64Type(), fmt.Sprintf("%d", decl.Id), 10), read_labels[i])
	}

	// XXX:
	//     event->XXX = minipack_unpack_XXX(ptr, sz);
	//     ptr += sz;
	//     if(sz != 0) goto loop else goto loop_skip;
	for i, decl := range read_decls {
		m.builder.SetInsertPointAtEnd(read_labels[i])

		m.printf("read.7\n")
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

		field := m.structgep(m.load(event, ""), decl.Index(), "")
		m.store(m.builder.CreateCall(m.module.NamedFunction(minipack_func_name), []llvm.Value{m.load(ptr, ""), sz}, ""), field)
		m.store(m.builder.CreateGEP(m.load(ptr, ""), []llvm.Value{m.load(sz, "")}, ""), ptr)
		m.printf("read.7.1 %d | (next=%p)\n", m.load(field), m.load(ptr))
		m.condbr(m.builder.CreateICmp(llvm.IntNE, m.load(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop, loop_skip)
	}

	// loop_skip:
	//     sz = minipack_sizeof_elem_and_data(ptr);
	//     ptr += sz;
	//     if(sz != 0) goto loop else goto error;
	m.builder.SetInsertPointAtEnd(loop_skip)
	m.store(m.builder.CreateCall(m.module.NamedFunction("minipack_sizeof_elem_and_data"), []llvm.Value{m.load(ptr, "")}, ""), sz)
	m.store(m.builder.CreateGEP(m.load(ptr, ""), []llvm.Value{m.load(sz, "")}, ""), ptr)
	m.printf("read.8\n")
	m.condbr(m.builder.CreateICmp(llvm.IntNE, m.load(sz, ""), llvm.ConstInt(m.context.Int64Type(), 0, false), ""), loop, error)

	// error:
	//     return false;
	m.builder.SetInsertPointAtEnd(error)
	m.printf("read.9\n")
	m.ret(llvm.ConstInt(m.context.Int1Type(), 0, false))

	// exit:
	//     return true;
	m.builder.SetInsertPointAtEnd(exit)
	m.printf("read.10\n")
	m.ret(llvm.ConstInt(m.context.Int1Type(), 1, false))

	return fn
}
