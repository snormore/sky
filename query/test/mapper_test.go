package test

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/axw/gollvm/llvm"
	"github.com/stretchr/testify/assert"
)

type event struct {
	eos bool
	eof bool
	timestamp int64
}

func TestReadEvent(t *testing.T) {
	query := `SELECT count()`
	m, _ := newMapper(query)
	e := new(event)
	m.ExecuteFunction("read_event", []llvm.GenericValue{
		llvm.NewGenericValueFromPointer(unsafe.Pointer(e)),
		llvm.NewGenericValueFromPointer(unsafe.Pointer(&[1]byte{'\x0a'})),
		})
	assert.Equal(t, int64(10), e.timestamp)
}
