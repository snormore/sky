package skyd

import (
  "fmt"
)

// Converts untyped map to a map[string]interface{} if passed a map.
func ConvertToStringKeys(value interface{}) interface{} {
  if m, ok := value.(map[interface{}]interface{}); ok {
    ret := make(map[string]interface{})
    for k,v := range m {
      ret[fmt.Sprintf("%v", k)] = v
    }
    return ret
  }
  
  return value
}
