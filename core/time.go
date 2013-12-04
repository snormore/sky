package core

import (
	"time"

	"github.com/skydb/sky/endian"
)

const SECONDS_BIT_OFFSET = 20

// Shifts a Go time into Sky timestamp format.
func ShiftTime(value time.Time) int64 {
	timestamp := value.UnixNano() / 1000
	usec := timestamp % 1000000
	sec := timestamp / 1000000
	return (sec << SECONDS_BIT_OFFSET) + usec
}

// ShiftTimeBytes converts a Go time in to a byte slice in Sky timestamp format.
func ShiftTimeBytes(value time.Time) []byte {
	var b [8]byte
	bs := b[:8]
	timestamp := ShiftTime(value)
	endian.Native.PutUint64(bs, uint64(timestamp))
	return bs
}

// Shifts a Sky timestamp format into a Go time.
func UnshiftTime(value int64) time.Time {
	usec := value & 0xFFFFF
	sec := value >> SECONDS_BIT_OFFSET
	return time.Unix(sec, usec*1000).UTC()
}

// UnshiftTimeBytes converts a byte slice in Sky timestamp format to Go time.
func UnshiftTimeBytes(value []byte) time.Time {
	timestamp := endian.Native.Uint64(value)
	return UnshiftTime(int64(timestamp))
}
