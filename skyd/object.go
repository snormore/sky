package skyd

import (
  "github.com/ugorji/go-msgpack"
)

// Encodes an object identifier.
func EncodeObjectId(tableName string, objectId string) ([]byte, error) {
  return msgpack.Marshal([]string{tableName,objectId})
}

