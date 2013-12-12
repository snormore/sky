package mapper

import (
	"fmt"
	
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
)

// [codegen]
// void field(sky_cursor *cursor, sky_map *result) {
// fields:
//     ...generate...
//
// exit:
//     return;
// }
func (m *Mapper) codegenField(node *ast.Field, tbl *ast.Symtable, index int) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.hashmapType, 0)}, false)
	fn := llvm.AddFunction(m.module, node.Identifier() + "_field", sig)

	entry := m.context.AddBasicBlock(fn, "entry")
	body := m.context.AddBasicBlock(fn, "body")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	result := m.alloca(llvm.PointerType(m.hashmapType, 0), "result")
	m.store(fn.Param(0), cursor)
	m.store(fn.Param(1), result)
	m.br(body)

	m.builder.SetInsertPointAtEnd(body)
	event := m.load(m.structgep(m.load(cursor), cursorEventElementIndex), "event")
	if node.IsAggregate() {
		currentValue := m.builder.CreateCall(m.module.NamedFunction("sky_hashmap_get"), []llvm.Value{m.load(result), llvm.ConstInt(m.context.Int64Type(), uint64(index), false)}, "")
		newValue, err := m.codegenAggregateField(node, tbl, event, currentValue)
		if err != nil {
			return nilValue, err
		}
		m.printf("field: %d -> %d\n", currentValue, newValue)
		m.builder.CreateCall(m.module.NamedFunction("sky_hashmap_set"), []llvm.Value{m.load(result), llvm.ConstInt(m.context.Int64Type(), uint64(index), false), newValue}, "")
	} else {
		panic("UNIMPLEMENTED")
	}
	m.br(exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.retvoid()

	return fn, nil
}

func (m *Mapper) codegenAggregateField(node *ast.Field, tbl *ast.Symtable, event llvm.Value, currentValue llvm.Value) (llvm.Value, error) {
	// The "count" aggregation is a special case since it doesn't require an expression.
	if node.Aggregation == "count" {
		return m.builder.CreateAdd(currentValue, llvm.ConstInt(m.context.Int64Type(), 1, false), ""), nil
	}
	
	// Generate the field expression.
	expressionValue, err := m.codegenExpression(node.Expression, event, tbl)
	if err != nil {
		return nilValue, err
	}

	// Aggregate the value of the expression with the current value.
	switch node.Aggregation {
	case "sum":
		return m.builder.CreateAdd(currentValue, expressionValue, ""), nil
	default:
		return nilValue, fmt.Errorf("Invalid aggregation type: %s", node.Aggregation)
	}
}

