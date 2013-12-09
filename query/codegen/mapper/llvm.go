package mapper

import (
	"github.com/axw/gollvm/llvm"
)

func (m *Mapper) alloca(typ llvm.Type, name... string) llvm.Value {
	return m.builder.CreateAlloca(typ, fname(name))
}

func (m *Mapper) br(bb llvm.BasicBlock) llvm.Value {
	return m.builder.CreateBr(bb)
}

func (m *Mapper) condbr(ifv llvm.Value, thenb, elseb llvm.BasicBlock) llvm.Value {
	return m.builder.CreateCondBr(ifv, thenb, elseb)
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
