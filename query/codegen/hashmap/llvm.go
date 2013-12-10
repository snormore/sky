package hashmap

import (
	"github.com/axw/gollvm/llvm"
)

func DeclareType(m llvm.Module, c llvm.Context) llvm.Type {
    return c.StructCreateNamed("sky_hashmap")
}

func Declare(m llvm.Module, c llvm.Context) {
    llvm.AddFunction(m, "sky_hashmap_new", llvm.FunctionType(c.Int64Type(), []llvm.Type{llvm.PointerType(c.Int8Type(), 0), llvm.PointerType(c.Int64Type(), 0)}, false))
}
