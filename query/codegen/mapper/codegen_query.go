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
	typ := m.context.StructCreateNamed("event_t")
	typ.StructSetBody(fields, false)
	return typ, nil
}

// [codegen]
// void entry(int32_t) {
//     return 12;
// }
func (m *Mapper) codegenQueryEntryFunc(q *ast.Query, tbl *symtable.Symtable) (llvm.Value, error) {
	fn := llvm.AddFunction(m.module, "entry", llvm.FunctionType(m.context.Int32Type(), []llvm.Type{m.eventType}, false))
	fn.SetFunctionCallConv(llvm.CCallConv)
	entry := m.context.AddBasicBlock(fn, "entry")
	m.builder.SetInsertPointAtEnd(entry)
	m.builder.CreateRet(llvm.ConstInt(m.context.Int32Type(), 12, false))
	return fn, nil
}

