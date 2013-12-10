package mapper

import (
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

// [codegen]
// void selection(sky_cursor *cursor, sky_map *result) {
// exit:
//     return;
// }
func (m *Mapper) codegenSelection(node *ast.Selection, tbl *ast.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.mapType, 0)}, false)
	fn := llvm.AddFunction(m.module, "selection", sig)

	// Generate functions for fields.
	var fieldFns []llvm.Value
	for _, field := range node.Fields {
		fieldFn, err := m.codegenField(field, tbl)
		if err != nil {
			return nilValue, err
		}
		fieldFns = append(fieldFns, fieldFn)
	}

	entry := m.context.AddBasicBlock(fn, "entry")
	fields_lbl := m.context.AddBasicBlock(fn, "fields")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor := m.builder.CreateAlloca(llvm.PointerType(m.cursorType, 0), "cursor")
	m.builder.CreateStore(fn.Param(0), cursor)
	result := m.builder.CreateAlloca(llvm.PointerType(m.mapType, 0), "result")
	m.builder.CreateStore(fn.Param(1), result)
	m.printf("selection.1\n")
	m.builder.CreateBr(fields_lbl)

	m.builder.SetInsertPointAtEnd(fields_lbl)
	for _, fieldFn := range fieldFns {
		m.builder.CreateCall(fieldFn, []llvm.Value{m.load(cursor, ""), m.load(result, "")}, "")
	}
	m.builder.CreateBr(exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.builder.CreateRetVoid()

	return fn, nil
}
