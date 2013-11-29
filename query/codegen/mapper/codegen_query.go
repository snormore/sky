package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/symtable"
)

func (m *Mapper) codegenQuery(q *ast.Query, tbl *symtable.Symtable) (llvm.Value, error) {
	// Generate "event" struct type.
	decls, err := ast.FindVarDecls(q)
	if err != nil {
		return nilValue, err
	}
	if m.eventType, err = m.codegenEventType(decls); err != nil {
		return nilValue, err
	}
	m.cursorType = m.context.StructCreateNamed("sky_cursor")
	m.mapType = m.context.StructCreateNamed("sky_map")

	m.cursorInitFunc = m.codegenCursorInitFunc()

	// Generate the entry function.
	return m.codegenQueryEntryFunc(q, tbl)
}

// [codegen]
// typedef struct {
//     ...
//     struct  { void *data; int64_t sz; } string_var;
//     int64_t factor_var;
//     int64_t integer_var;
//     double  float_var;
//     bool    boolean_var;
//     ...
// } event_t;
func (m *Mapper) codegenEventType(decls ast.VarDecls) (llvm.Type, error) {
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
	return typ, nil
}

// [codegen]
// bool sky_cursor_init(sky_cursor *)
func (m *Mapper) codegenCursorInitFunc() llvm.Value {
	return llvm.AddFunction(m.module, "sky_cursor_init", llvm.FunctionType(m.context.Int1Type(), []llvm.Type{llvm.PointerType(m.cursorType, 0)}, false))
}

// [codegen]
// int32_t entry(sky_cursor *cursor, sky_map *result) {
//     int32_t rc = 0;
//     entry:
//         rc = sky_cursor_init();
//         if(rc != EEOF) {
//             goto exit;
//         }
//         if(rc != 0) {
//             goto error;
//         }
//
//     loop:
//         rc = sky_cursor_next_object();
//         if(rc != EEOF) {
//             goto exit;
//         }
//         if(rc != 0) {
//             goto error;
//         }
//         goto loop
//
//     error:
//     exit:
//         return rc;
// }
func (m *Mapper) codegenQueryEntryFunc(q *ast.Query, tbl *symtable.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.Int32Type(), []llvm.Type{
		llvm.PointerType(m.cursorType, 0),
		llvm.PointerType(m.mapType, 0),
	}, false)
	fn := llvm.AddFunction(m.module, "entry", sig)
	fn.SetFunctionCallConv(llvm.CCallConv)
	cursorArg := fn.Param(0)
	// resultArg := fn.Param(1)

	entry := m.context.AddBasicBlock(fn, "entry")
	//loop := m.context.AddBasicBlock(fn, "loop")
	//error := m.context.AddBasicBlock(fn, "error")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	m.builder.CreateCall(m.cursorInitFunc, []llvm.Value{cursorArg}, "")
	m.builder.CreateBr(exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateRet(llvm.ConstInt(m.context.Int32Type(), 12, false))
	return fn, nil
}
