package skyd

import (
  "testing"
)

// Ensure that the lua script can initialize.
func TestLuaScriptInit(t *testing.T) {
  l := NewLuaScript("print('hello')")
  l.ExecuteFunction("test")
  l.Destroy()
}
