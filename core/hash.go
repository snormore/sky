package core

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"strings"
)

//------------------------------------------------------------------------------
//
// Functions
//
//------------------------------------------------------------------------------

//--------------------------------------
// Id Generation
//--------------------------------------

// Generates a new short identifier. This is not guarenteed to be unique. It's
// just random.
func NewId(length int) string {
	s := fmt.Sprintf("%x", rand.Int63())
	if len(s) < length {
		s = strings.Repeat("0", length-len(s)) + s
	}
	return s[0:length]
}

//--------------------------------------
// Object Hash
//--------------------------------------

// Generates the FNV1a hash for an object identifier.
func ObjectHash(objectId string) uint64 {
	h := fnv.New64a()
	h.Reset()
	h.Write([]byte(objectId))
	return h.Sum64()
}

// Retrieves the core-level hash for an object.
func ObjectLocalHash(objectId string) uint32 {
	return uint32(ObjectHash(objectId) & 0xFFFFFFFF)
}

// Retrieves the node-level hash for an object.
func ObjectShardHash(objectId string) uint32 {
	return uint32(ObjectHash(objectId) >> 32)
}
