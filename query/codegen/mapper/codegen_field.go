package mapper

import (
	"fmt"
	
	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/codegen/hash"
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
		id := hash.HashString(node.Identifier())
		currentValue := m.call("sky_hashmap_get", m.load(result), m.constint(int(id)))
		newValue, err := m.codegenAggregateField(node, tbl, event, currentValue)
		if err != nil {
			return nilValue, err
		}
		m.printf("field.set %s %d -> %d\n", node.Identifier(), currentValue, newValue)
		m.call("sky_hashmap_set", m.load(result), m.constint(int(id)), newValue)
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

