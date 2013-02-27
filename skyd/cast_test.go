package skyd

import (
	"testing"
)

// Cast integer values to Int64.
func TestCastInt64(t *testing.T) {
  var i int = 100
  var i8 int8 = 100
  var i16 int16 = 100
  var i32 int32 = 100
  var i64 int64 = 100

  var u int = 100
  var u8 uint8 = 100
  var u16 uint16 = 100
  var u32 uint32 = 100
  var u64 uint64 = 100
  
  var ret int64
  
  ret, _ = castInt64(i)
  if ret != 100 {
    t.Errorf("Unable to cast int")
  }
  ret, _ = castInt64(i8)
  if ret != 100 {
    t.Errorf("Unable to cast int8")
  }
  ret, _ = castInt64(i16)
  if ret != 100 {
    t.Errorf("Unable to cast int16")
  }
  ret, _ = castInt64(i32)
  if ret != 100 {
    t.Errorf("Unable to cast int32")
  }
  ret, _ = castInt64(i64)
  if ret != 100 {
    t.Errorf("Unable to cast int64")
  }

  ret, _ = castInt64(u)
  if ret != 100 {
    t.Errorf("Unable to cast uint")
  }
  ret, _ = castInt64(u8)
  if ret != 100 {
    t.Errorf("Unable to cast uint8")
  }
  ret, _ = castInt64(u16)
  if ret != 100 {
    t.Errorf("Unable to cast uint16")
  }
  ret, _ = castInt64(u32)
  if ret != 100 {
    t.Errorf("Unable to cast uint32")
  }
  ret, _ = castInt64(u64)
  if ret != 100 {
    t.Errorf("Unable to cast uint64")
  }
}

