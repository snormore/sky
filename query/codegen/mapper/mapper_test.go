package mapper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/axw/gollvm/llvm"
	"github.com/skydb/sky/core"
	"github.com/skydb/sky/endian"
	"github.com/skydb/sky/query/ast"
	"github.com/skydb/sky/query/parser"
	"github.com/stretchr/testify/assert"
)

type event struct {
	eos bool
	eof bool
	timestamp int64
	foo int64
}

var DATA = [1024]byte{}

func TestReadEvent(t *testing.T) {
	var b bytes.Buffer
	binary.Write(&b, endian.Native, int64(10 << core.SECONDS_BIT_OFFSET))
	b.Write([]byte{'\x81', '\x02', '\x03'})  // map: {2:3}
	copy(DATA[:], b.Bytes())
	debugln(DATA)

	query, _ := parser.ParseString(`SELECT count()`)
	query.DeclaredVarDecls = append(query.DeclaredVarDecls, ast.NewVarDecl(2, "foo", "integer"))
	m, err := New(query)
	assert.NoError(t, err)
	e := new(event)
	m.Dump()
	m.ExecuteFunction("read_event", []llvm.GenericValue{
		llvm.NewGenericValueFromPointer(unsafe.Pointer(e)),
		llvm.NewGenericValueFromPointer(unsafe.Pointer(&DATA)),
		})
	assert.Equal(t, int64(10), e.timestamp)
	assert.Equal(t, int64(3), e.foo)
}


func debugln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
