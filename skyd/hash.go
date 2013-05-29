package skyd

import (
	"fmt"
	"math/rand"
	"strings"
)

//------------------------------------------------------------------------------
//
// Functions
//
//------------------------------------------------------------------------------

// Generates a new short identifier. This is not guarenteed to be unique. It's
// just random.
func NewId(length int) string {
	s := fmt.Sprintf("%x", rand.Int63())
	if len(s) < length {
		s = strings.Repeat("0", length-len(s)) + s
	}
	return s[0:length]
}
