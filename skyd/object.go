package skyd

import (
  "github.com/ugorji/go-msgpack"
)

// Encodes an object identifier.
func EncodeObjectId(objectId interface{}) ([]byte, error) {
  return msgpack.Marshal(objectId)
}

