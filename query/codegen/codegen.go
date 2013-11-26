package codegen

import (
	"github.com/axw/gollvm/llvm"
)

func init() {
	llvm.LinkInJIT()
	llvm.InitializeNativeTarget()
}
