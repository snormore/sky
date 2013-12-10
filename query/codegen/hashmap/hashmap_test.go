package hashmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Ensure that integer values can be set and retrieved.
func TestHashmapInt(t *testing.T) {
	h := New()
	defer h.Free()

	max := 10000
	for i := 0; i < max; i++ {
		h.Set(int64(i), int64(max - i))
	}
	for i := 0; i < max; i++ {
		assert.Equal(t, h.Get(int64(i)), int64(max - i))
	}
}

// Ensure that hashmap values can be retrieved.
func TestHashmapSubmap(t *testing.T) {
	h := New()
	defer h.Free()
	s0 := h.Submap(10)
	s0.Set(100, 200)
	s1 := h.Submap(11)
	s1.Set(100, 300)
	assert.Equal(t, h.Submap(10).Get(100), 200)
	assert.Equal(t, h.Submap(11).Get(100), 300)
}

// Ensure that retrieving a submap will overwrite an int value.
func TestHashmapOverrideWithSubmap(t *testing.T) {
	h := New()
	defer h.Free()
	h.Set(10, 100)
	s0 := h.Submap(10)
	s0.Set(20, 200)
	assert.Equal(t, h.Submap(10).Get(20), 200)
}

// Ensure that retrieving an int will overwrite a submap value.
func TestHashmapOverrideWithInt(t *testing.T) {
	h := New()
	defer h.Free()
	h.Submap(10)
	h.Set(10, 100)
	assert.Equal(t, h.Get(10), 100)
}


func BenchmarkHashmapInsert(b *testing.B) {
	h := New()
	for i := 0; i < b.N; i++ {
		key := i % 1024
		h.Set(int64(key), 100)
	}
}
