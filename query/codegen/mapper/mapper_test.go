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
	eos          bool
	eof          bool
	timestamp    int64
	booleanValue bool
	integerValue int64
}

var DATA = [1024]byte{}

var mapperVarDecls = ast.VarDecls{
	ast.NewVarDecl(2, "integerValue", "integer"),
	ast.NewVarDecl(-1, "booleanValue", "boolean"),
}

func TestReadEvent(t *testing.T) {
	var b bytes.Buffer
	binary.Write(&b, endian.Native, int64(10<<core.SECONDS_BIT_OFFSET))
	b.Write([]byte{'\x82', '\x02', '\x03', '\xFF', '\xC3'}) // map: {2:3, -1:true}
	copy(DATA[:], b.Bytes())

	query, _ := parser.ParseString(`SELECT count()`)
	query.DeclaredVarDecls = append(query.DeclaredVarDecls, mapperVarDecls...)
	m, err := New(query)
	if assert.NoError(t, err) {
		e := new(event)
		m.ExecuteFunction("read_event", []llvm.GenericValue{
			llvm.NewGenericValueFromPointer(unsafe.Pointer(e)),
			llvm.NewGenericValueFromPointer(unsafe.Pointer(&DATA)),
		})
		assert.Equal(t, int64(10), e.timestamp)
		assert.Equal(t, int64(3), e.integerValue)
		assert.Equal(t, true, e.booleanValue)
	}
}

func debugln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
