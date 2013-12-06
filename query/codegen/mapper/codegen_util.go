package mapper

import (
	"github.com/axw/gollvm/llvm"
)

func (m *Mapper) createPrintfCall(format string, values ...interface{}) llvm.Value {
	vals := []llvm.Value{}
	vals = append(vals, m.builder.CreateGlobalString(format, ""))

	for _, value := range values {
		switch value := value.(type) {
		case string:
			vals = append(vals, m.builder.CreateGlobalString(value, ""))
		case llvm.Value:
			vals = append(vals, value)
		default:
			panic("Invalid argument to printf call!")
		}
	}
	return m.builder.CreateCall(m.module.NamedFunction("printf"), vals, "")
}
