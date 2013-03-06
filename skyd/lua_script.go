package skyd

/*
#cgo LDFLAGS: -lluajit-5.1.2
#include <stdlib.h>
#include <lua.h>
#include <lualib.h>
#include <lauxlib.h>
*/
import "C"

import (
  "errors"
  "fmt"
  "unsafe"
)

// A LuaScript is used to bridge between Go and LuaJIT.
type LuaScript struct {
  state   *C.lua_State
  source  string
}

func NewLuaScript(source string) *LuaScript {
  return &LuaScript{source:source}
}

// Executes the lua script and returns the data.
func (l *LuaScript) ExecuteFunction(functionName string, v ...interface{}) (interface{}, error) {
  if err := l.Init(); err != nil {
    return nil, err
  }

  // TODO: Execute function.
  return nil, nil
}

// Initializes the Lua context and compiles the source code.
func (l *LuaScript) Init() error {
  if l.state != nil {
    return nil
  }

  // Initialize the state and open the libraries.
  l.state = C.luaL_newstate()
  if l.state == nil {
    return errors.New("Unable to initialize Lua context.")
  }
  C.luaL_openlibs(l.state);
  
  // Compile the script.
  source := C.CString(l.source)
  defer C.free(unsafe.Pointer(source))
  ret := C.luaL_loadstring(l.state, source)
  if ret != 0 {
    errstring := C.GoString(C.lua_tolstring(l.state, -1, nil))
    return fmt.Errorf("skyd.LuaScript: Syntax Error: %v", errstring)
  }

  // Run script once to initialize.
  ret = C.lua_pcall(l.state, 0, 0, 0);
  if ret != 0 {
    errstring := C.GoString(C.lua_tolstring(l.state, -1, nil))
    return fmt.Errorf("skyd.LuaScript: Init Error: %v", errstring)
  }
  
  return nil
}

// Closes the lua context.
func (l *LuaScript) Destroy() {
  if l.state != nil {
    C.lua_close(l.state)
  }
}