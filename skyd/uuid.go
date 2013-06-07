package skyd

import (
	"crypto/rand"
	"fmt"
)

type UUID [16]byte

func NewUUID() *UUID {
	u := new(UUID)
	rand.Read(u[:])
	return u
}

func (u *UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
}
