package mapper

import (
	"fmt"

	"github.com/axw/gollvm/llvm"
)

var nilValue llvm.Value

func (m *Mapper) alloca(typ llvm.Type, name... string) llvm.Value {
	return m.builder.CreateAlloca(typ, fname(name))
}

func (m *Mapper) add(lhs llvm.Value, rhs llvm.Value, name...string) llvm.Value {
	return m.builder.CreateAdd(lhs, rhs, fname(name))
}

func (m *Mapper) br(bb llvm.BasicBlock) llvm.Value {
	return m.builder.CreateBr(bb)
}

func (m *Mapper) call(fn interface{}, args... llvm.Value) llvm.Value {
	if functionName, ok := fn.(string); ok {
		fn = m.module.NamedFunction(functionName)
	}
	return m.builder.CreateCall(fn.(llvm.Value), args, "")
}

func (m *Mapper) condbr(ifv llvm.Value, thenb, elseb llvm.BasicBlock) llvm.Value {
	return m.builder.CreateCondBr(ifv, thenb, elseb)
}

func (m *Mapper) constbool(value bool) llvm.Value {
	if value {
		return llvm.ConstInt(m.context.Int1Type(), 1, false)
	}
	return llvm.ConstInt(m.context.Int1Type(), 0, false)
}

func (m *Mapper) constint(value int) llvm.Value {
	return llvm.ConstIntFromString(m.context.Int64Type(), fmt.Sprintf("%d", value), 10)
}

func (m *Mapper) icmp(pred llvm.IntPredicate, lhs, rhs llvm.Value, name... string) llvm.Value {
	return m.builder.CreateICmp(pred, lhs, rhs, fname(name))
}

func (m *Mapper) load(value llvm.Value, name... string) llvm.Value {
	return m.builder.CreateLoad(value, fname(name))
}

func (m *Mapper) ret(value llvm.Value) llvm.Value {
	return m.builder.CreateRet(value)
}

func (m *Mapper) retvoid() llvm.Value {
	return m.builder.CreateRetVoid()
}

func (m *Mapper) store(src llvm.Value, dest llvm.Value) llvm.Value {
	return m.builder.CreateStore(src, dest)
}

func (m *Mapper) structgep(value llvm.Value, index int, name... string) llvm.Value {
	return m.builder.CreateStructGEP(value, index, fname(name))
}


func fname(names []string) string {
	if len(names) > 0 {
		return names[0]
	}
	return ""
}
