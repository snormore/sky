package core

import (
	"math/rand"
	"testing"
)

// Ensure that a identifiers are generated correctly.
func TestHashNewId(t *testing.T) {
	rand.Seed(1)
	if id := NewId(7); id != "4d65822" {
		t.Fatalf("Unexpected identifier: %v", id)
	}
	if id := NewId(20); id != "000078629a0f5f3f164f" {
		t.Fatalf("Unexpected identifier: %v", id)
	}
}

// Ensure that a identifiers are padded correctly.
func TestHashNewIdWithPadding(t *testing.T) {
	rand.Seed(1)
	if id := NewId(20); id != "00004d65822107fcfd52" {
		t.Fatalf("Unexpected identifier: %v", id)
	}
}

func BenchmarkHashNewId(b *testing.B) {
	rand.Seed(1)
	for i := 0; i < b.N; i++ {
		NewId(7)
	}
}
