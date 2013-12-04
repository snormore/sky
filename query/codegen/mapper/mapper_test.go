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
	"github.com/skydb/sky/query/parser"
	"github.com/stretchr/testify/assert"
)

type event struct {
	eos bool
	eof bool
	timestamp int64
}

var DATA = [1024]byte{}

func TestReadEvent(t *testing.T) {
	var b bytes.Buffer
	binary.Write(&b, endian.Native, int64(10 << core.SECONDS_BIT_OFFSET))
	b.Write([]byte{'\x81', '\x02', '\x03'})  // map: {2:3}
	copy(DATA[:], b.Bytes())
	debugln(DATA)

	query := `SELECT count()`
	m, err := newMapperFromQuery(query)
	assert.NoError(t, err)
	e := new(event)
	m.Dump()
	m.ExecuteFunction("read_event", []llvm.GenericValue{
		llvm.NewGenericValueFromPointer(unsafe.Pointer(e)),
		llvm.NewGenericValueFromPointer(unsafe.Pointer(&DATA)),
		})
	assert.Equal(t, int64(10), e.timestamp)
}


func newMapperFromQuery(query string) (*Mapper, error) {
	// Create a query.
	q, err := parser.ParseString(query)
	if err != nil {
		return nil, err
	}

	// Create a mapper generated from the query.
	m, err := New(q)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func debugln(a ...interface{}) (n int, err error) {
	return fmt.Fprintln(os.Stderr, a...)
}
