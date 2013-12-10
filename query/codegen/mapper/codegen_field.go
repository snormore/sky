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
func (m *Mapper) codegenField(node *ast.Field, tbl *ast.Symtable) (llvm.Value, error) {
	sig := llvm.FunctionType(m.context.VoidType(), []llvm.Type{llvm.PointerType(m.cursorType, 0), llvm.PointerType(m.mapType, 0)}, false)
	fn := llvm.AddFunction(m.module, node.Identifier() + "_field", sig)

	entry := m.context.AddBasicBlock(fn, "entry")
	body := m.context.AddBasicBlock(fn, "body")
	exit := m.context.AddBasicBlock(fn, "exit")

	m.builder.SetInsertPointAtEnd(entry)
	cursor := m.alloca(llvm.PointerType(m.cursorType, 0), "cursor")
	result := m.alloca(llvm.PointerType(m.mapType, 0), "result")
	currentValue := m.alloca(m.context.Int64Type(), "currentValue")
	m.store(fn.Param(0), cursor)
	m.store(fn.Param(1), result)
	m.store(llvm.ConstInt(m.context.Int64Type(), 0, false), currentValue)
	m.br(body)

	m.builder.SetInsertPointAtEnd(body)
	event := m.load(m.structgep(m.load(cursor), cursorEventElementIndex), "event")
	if node.IsAggregate() {
		if _, err := m.codegenAggregateField(node, event, m.load(currentValue), tbl); err != nil {
			return nilValue, err
		}
	} else {
		if _, err := m.codegenNonAggregateField(node, event, tbl); err != nil {
			return nilValue, err
		}
	}
	m.br(exit)

	m.builder.SetInsertPointAtEnd(exit)
	m.retvoid()

	return fn, nil
}

func (m *Mapper) codegenAggregateField(node *ast.Field, event llvm.Value, currentValue llvm.Value, tbl *ast.Symtable) (llvm.Value, error) {
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

func (m *Mapper) codegenNonAggregateField(node *ast.Field, event llvm.Value, tbl *ast.Symtable) (llvm.Value, error) {
	panic("UNIMPLEMENTED")
}

